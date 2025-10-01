package router

import (
	"strings"
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
	n.priority++

	// Empty path means this node is the target
	if path == "" {
		if n.methods == nil {
			n.methods = make(map[string]HandlerFunc[T])
		}
		n.methods[method] = handler
		return
	}

	// Safety check to prevent infinite recursion
	if len(path) > 1000 {
		panic("path too long, possible infinite recursion")
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
	n.insertStaticRoute(method, path, handler)
}

// insertStaticRoute handles static path segments
func (n *node[T]) insertStaticRoute(method, path string, handler HandlerFunc[T]) {
	// Find the first slash or end of string
	slashIndex := strings.Index(path, "/")
	var staticPart string
	var remainingPath string

	if slashIndex == -1 {
		// No slash found, entire path is static
		staticPart = path
		remainingPath = ""
	} else {
		// Split at the slash
		staticPart = path[:slashIndex]
		remainingPath = path[slashIndex+1:]
	}

	// If this node has no path yet, set it
	if n.path == "" {
		n.path = staticPart
		n.nType = static
		if remainingPath == "" {
			// This is the final node
			if n.methods == nil {
				n.methods = make(map[string]HandlerFunc[T])
			}
			n.methods[method] = handler
		} else {
			// Continue with remaining path
			n.addRoute(method, remainingPath, handler)
		}
		return
	}

	// Find common prefix with current node
	commonLen := 0
	maxLen := min(len(staticPart), len(n.path))
	for commonLen < maxLen && staticPart[commonLen] == n.path[commonLen] {
		commonLen++
	}

	// If we have a common prefix, we need to split this node
	if commonLen < len(n.path) {
		// Split the current node
		child := &node[T]{
			path:       n.path[commonLen:],
			nType:      n.nType,
			children:   n.children,
			methods:    n.methods,
			indices:    n.indices,
			wildChild:  n.wildChild,
			paramChild: n.paramChild,
			priority:   n.priority - 1,
		}

		// Reset current node
		n.path = n.path[:commonLen]
		n.children = []*node[T]{child}
		n.indices = []byte{child.path[0]}
		n.methods = nil
		n.wildChild = nil
		n.paramChild = nil
	}

	// If we've consumed the entire static part, continue with remaining path
	if commonLen == len(staticPart) {
		if remainingPath == "" {
			// This is the final node
			if n.methods == nil {
				n.methods = make(map[string]HandlerFunc[T])
			}
			n.methods[method] = handler
		} else {
			// Continue with remaining path
			n.addRoute(method, remainingPath, handler)
		}
		return
	}

	// We need to add a new child for the remaining static part
	remainingStatic := staticPart[commonLen:]

	if len(remainingStatic) == 0 {
		// This shouldn't happen, but handle it gracefully
		if remainingPath == "" {
			if n.methods == nil {
				n.methods = make(map[string]HandlerFunc[T])
			}
			n.methods[method] = handler
		} else {
			n.addRoute(method, remainingPath, handler)
		}
		return
	}

	c := remainingStatic[0]

	// Check if we already have a child with this character
	for i, index := range n.indices {
		if index == c {
			n.children[i].addRoute(method, remainingStatic+"/"+remainingPath, handler)
			return
		}
	}

	// Create new child
	child := &node[T]{
		nType:   static,
		methods: make(map[string]HandlerFunc[T]),
	}

	// Add the child
	n.addChild(child, c)

	// Set up the child's path and continue
	if remainingPath == "" {
		child.path = remainingStatic
		child.methods[method] = handler
	} else {
		child.path = remainingStatic
		child.addRoute(method, remainingPath, handler)
	}
}

// insertParamRoute handles parameter routes (:param)
func (n *node[T]) insertParamRoute(method, path string, handler HandlerFunc[T]) {
	// Find parameter name (until next slash or end)
	end := 1
	for end < len(path) && path[end] != '/' {
		end++
	}

	paramName := path[1:end]

	// Create or get parameter child
	if n.paramChild == nil {
		n.paramChild = &node[T]{
			nType:     param,
			paramName: paramName,
			methods:   make(map[string]HandlerFunc[T]),
		}
	}

	// Continue with remaining path
	if end < len(path) {
		n.paramChild.addRoute(method, path[end+1:], handler)
	} else {
		if n.paramChild.methods == nil {
			n.paramChild.methods = make(map[string]HandlerFunc[T])
		}
		n.paramChild.methods[method] = handler
	}
}

// insertWildcardRoute handles wildcard routes (*wildcard)
func (n *node[T]) insertWildcardRoute(method, path string, handler HandlerFunc[T]) {
	// Find wildcard name (until next slash or end)
	end := 1
	for end < len(path) && path[end] != '/' {
		end++
	}

	paramName := path[1:end]

	// Create or get wildcard child
	if n.wildChild == nil {
		n.wildChild = &node[T]{
			nType:     wildcard,
			paramName: paramName,
			methods:   make(map[string]HandlerFunc[T]),
		}
	}

	// Wildcard consumes rest of path
	if n.wildChild.methods == nil {
		n.wildChild.methods = make(map[string]HandlerFunc[T])
	}
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
	// If we have a path, check if it matches
	if n.path != "" {
		if len(path) < len(n.path) || path[:len(n.path)] != n.path {
			return nil
		}
		path = path[len(n.path):]
		// If there's a slash after the matched path, consume it
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}
	}

	// If we've consumed all path, check for handler
	if path == "" {
		if handler, ok := n.methods[method]; ok {
			return handler
		}
		return nil
	}

	// Try static children first (highest priority)
	if len(n.children) > 0 {
		c := path[0]
		for i, index := range n.indices {
			if index == c {
				if handler := n.children[i].findRoute(method, path, params); handler != nil {
					return handler
				}
				break
			}
		}
	}

	// Try parameter child (medium priority)
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
		} else {
			// Continue with remaining path
			if handler := n.paramChild.findRoute(method, path[end+1:], params); handler != nil {
				return handler
			}
		}
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

	return nil
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
