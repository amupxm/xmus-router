package router

import (
	"sync"
	"unsafe"
)

// Parameter represents a URL parameter
type Parameter struct {
	Key   string
	Value string
}

// Parameters is a slice of Parameter with optimized access
type Parameters []Parameter

// Get retrieves a parameter value by key (zero-allocation lookup)
func (ps Parameters) Get(key string) (string, bool) {
	for _, p := range ps {
		if p.Key == key {
			return p.Value, true
		}
	}
	return "", false
}

// MustGet retrieves a parameter value by key, panics if not found
func (ps Parameters) MustGet(key string) string {
	if value, ok := ps.Get(key); ok {
		return value
	}
	panic("parameter not found: " + key)
}

// nodeType represents different node types for optimization
type nodeType uint8

const (
	static   nodeType = iota // static path segment
	param                    // :param
	wildcard                 // *wildcard
)

// node represents a radix tree node with generic context
type node[T Context] struct {
	// Path segment for this node
	path string

	// Node type for fast matching
	nType nodeType

	// Parameter name (for param/wildcard nodes)
	paramName string

	// Handler for this exact path (if any)
	handler Handler[T]

	// HTTP methods -> handlers mapping for this path
	methods map[string]HandlerFunc[T]

	// Children nodes
	children []*node[T]

	// Wildcard child (for * routes)
	wildChild *node[T]

	// Parameter child (for : routes)
	paramChild *node[T]

	// Indices for fast child lookup (first char of each child path)
	indices []byte

	// Priority for reordering (most used routes first)
	priority uint32
}

// radixTree represents the main router tree
type radixTree[T Context] struct {
	root *node[T]
	mu   sync.RWMutex // Thread safety
}

// NewRadixTree creates a new radix tree router
func NewRadixTree[T Context]() *radixTree[T] {
	return &radixTree[T]{
		root: &node[T]{
			methods: make(map[string]HandlerFunc[T]),
		},
	}
}

// Add inserts a new route into the tree
func (t *radixTree[T]) Add(method, path string, handler HandlerFunc[T]) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if path == "" || path[0] != '/' {
		panic("path must begin with '/'")
	}

	t.root.addRoute(method, path[1:], handler) // Remove leading /
	t.root.updatePriority()
}

// Find searches for a route and returns handler + parameters
func (t *radixTree[T]) Find(method, path string) (HandlerFunc[T], Parameters) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if path == "" || path[0] != '/' {
		return nil, nil
	}

	// Use pre-allocated slice to avoid allocations
	params := make(Parameters, 0, 8)
	handler := t.root.findRoute(method, path[1:], &params)

	return handler, params
}

// addRoute adds a route to the node
func (n *node[T]) addRoute(method, path string, handler HandlerFunc[T]) {
	fullPath := path
	n.priority++

	// Empty path means this node is the target
	if path == "" {
		if n.methods == nil {
			n.methods = make(map[string]HandlerFunc[T])
		}
		n.methods[method] = handler
		return
	}

	// Handle parameter routes (:param)
	if path[0] == ':' {
		n.insertParamRoute(method, path, handler)
		return
	}

	// Handle wildcard routes (*wildcard)
	if path[0] == '*' {
		n.insertWildcardRoute(method, path, handler)
		return
	}

	// Handle static routes
	n.insertStaticRoute(method, path, fullPath, handler)
}

// insertStaticRoute handles static path segments
func (n *node[T]) insertStaticRoute(method, path, fullPath string, handler HandlerFunc[T]) {
	// Find common prefix
	i := 0
	max := min(len(path), len(n.path))
	for i < max && path[i] == n.path[i] {
		i++
	}

	// Split node if needed
	if i < len(n.path) {
		child := &node[T]{
			path:       n.path[i:],
			nType:      n.nType,
			children:   n.children,
			methods:    n.methods,
			indices:    n.indices,
			wildChild:  n.wildChild,
			paramChild: n.paramChild,
			priority:   n.priority - 1,
		}

		n.children = []*node[T]{child}
		n.indices = []byte{n.path[i]}
		n.path = path[:i]
		n.methods = nil
		n.wildChild = nil
		n.paramChild = nil
	}

	// Add remaining path
	if i < len(path) {
		path = path[i:]
		c := path[0]

		// Find existing child
		for j, index := range n.indices {
			if c == index {
				n.children[j].addRoute(method, path, handler)
				return
			}
		}

		// Create new child
		child := &node[T]{
			nType:   static,
			methods: make(map[string]HandlerFunc[T]),
		}

		n.addChild(child, c)
		child.addRoute(method, path, handler)
	} else {
		// This node is the target
		if n.methods == nil {
			n.methods = make(map[string]HandlerFunc[T])
		}
		n.methods[method] = handler
	}
}

// insertParamRoute handles parameter routes (:param)
func (n *node[T]) insertParamRoute(method, path string, handler HandlerFunc[T]) {
	// Find parameter name
	end := 1
	for end < len(path) && path[end] != '/' {
		end++
	}

	paramName := path[1:end]

	if n.paramChild == nil {
		n.paramChild = &node[T]{
			nType:     param,
			paramName: paramName,
			methods:   make(map[string]HandlerFunc[T]),
		}
	}

	if end < len(path) {
		n.paramChild.addRoute(method, path[end+1:], handler)
	} else {
		n.paramChild.methods[method] = handler
	}
}

// insertWildcardRoute handles wildcard routes (*wildcard)
func (n *node[T]) insertWildcardRoute(method, path string, handler HandlerFunc[T]) {
	// Find wildcard name
	end := 1
	for end < len(path) && path[end] != '/' {
		end++
	}

	paramName := path[1:end]

	if n.wildChild == nil {
		n.wildChild = &node[T]{
			nType:     wildcard,
			paramName: paramName,
			methods:   make(map[string]HandlerFunc[T]),
		}
	}

	// Wildcard consumes rest of path
	n.wildChild.methods[method] = handler
}

// addChild adds a child node with proper ordering
func (n *node[T]) addChild(child *node[T], index byte) {
	// Insert in sorted order for faster lookup
	pos := 0
	for pos < len(n.indices) && n.indices[pos] < index {
		pos++
	}

	// Insert at position
	n.indices = append(n.indices, 0)
	copy(n.indices[pos+1:], n.indices[pos:])
	n.indices[pos] = index

	n.children = append(n.children, nil)
	copy(n.children[pos+1:], n.children[pos:])
	n.children[pos] = child
}

// findRoute searches for a route in the tree
func (n *node[T]) findRoute(method, path string, params *Parameters) HandlerFunc[T] {
walk:
	for {
		// Check if we've consumed all path
		if len(path) <= len(n.path) {
			if path == n.path {
				if handler, ok := n.methods[method]; ok {
					return handler
				}
			}
			return nil
		}

		// Check path prefix
		if path[:len(n.path)] == n.path {
			path = path[len(n.path):]

			// Try static children first (fastest)
			if len(n.indices) > 0 {
				c := path[0]

				// Binary search for better performance with many children
				i := 0
				j := len(n.indices)
				for i < j {
					mid := (i + j) / 2
					if n.indices[mid] < c {
						i = mid + 1
					} else {
						j = mid
					}
				}

				if i < len(n.indices) && n.indices[i] == c {
					n = n.children[i]
					continue walk
				}
			}

			// Try parameter child
			if n.paramChild != nil {
				// Find end of parameter value
				end := 0
				for end < len(path) && path[end] != '/' {
					end++
				}

				// Add parameter
				*params = append(*params, Parameter{
					Key:   n.paramChild.paramName,
					Value: path[:end],
				})

				if end == len(path) {
					// End of path, check for handler
					if handler, ok := n.paramChild.methods[method]; ok {
						return handler
					}
					return nil
				}

				// Continue with remaining path
				n = n.paramChild
				path = path[end+1:] // Skip the /
				continue walk
			}

			// Try wildcard child (lowest priority)
			if n.wildChild != nil {
				*params = append(*params, Parameter{
					Key:   n.wildChild.paramName,
					Value: path,
				})

				if handler, ok := n.wildChild.methods[method]; ok {
					return handler
				}
			}
		}

		return nil
	}
}

// updatePriority reorders children based on priority
func (n *node[T]) updatePriority() {
	// Sort children by priority (descending)
	for i := 1; i < len(n.children); i++ {
		child := n.children[i]
		index := n.indices[i]

		j := i
		for j > 0 && n.children[j-1].priority < child.priority {
			n.children[j] = n.children[j-1]
			n.indices[j] = n.indices[j-1]
			j--
		}

		n.children[j] = child
		n.indices[j] = index
	}

	// Recursively update children
	for _, child := range n.children {
		child.updatePriority()
	}

	if n.paramChild != nil {
		n.paramChild.updatePriority()
	}

	if n.wildChild != nil {
		n.wildChild.updatePriority()
	}
}

// Helper function for minimum
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Performance optimizations using unsafe for zero-allocation string operations
// Use these carefully and only when performance is critical

// unsafeString converts byte slice to string without allocation
func unsafeString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// unsafeBytes converts string to byte slice without allocation
func unsafeBytes(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	return *(*[]byte)(unsafe.Pointer(&s))
}
