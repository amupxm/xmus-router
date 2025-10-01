package router

import (
	"net/http"
)

// RouterOptions contains configuration for the router
type RouterOptions struct {
	NotFoundHandler  HandlerFunc[Context]
	MethodNotAllowed HandlerFunc[Context]
	CustomPrintf     func(format string, args ...any)
}

// Router wraps the radix tree with additional functionality
type Router struct {
	tree       *radixTree[Context]
	middleware []Middleware[Context]
	options    *RouterOptions
	groups     []*Group
}

// Group represents a route group with middleware
type Group struct {
	router     *Router
	prefix     string
	middleware []Middleware[Context]
	parent     *Group
}

// NewRouter creates a new high-performance router
func NewRouter(options *RouterOptions) *Router {
	if options == nil {
		options = &RouterOptions{}
	}
	return &Router{
		tree:       NewRadixTree[Context](),
		middleware: []Middleware[Context]{},
		options:    options,
		groups:     []*Group{},
	}
}

type HandlerFunc[T Context] func(w http.ResponseWriter, r *http.Request, ctx T)

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Create context
	ctx := NewContext(req, w)

	// Find route and parameters
	handler, params := r.tree.Find(req.Method, req.URL.Path)

	if handler == nil {
		// Try to find any handler for this path (for method not allowed)
		_, _ = r.tree.Find("", req.URL.Path)
		if r.options.MethodNotAllowed != nil {
			r.options.MethodNotAllowed(w, req, ctx)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Set parameters in context
	if len(params) > 0 {
		paramMap := make(map[string]string)
		for _, p := range params {
			paramMap[p.Key] = p.Value
		}
		ctx.SetParams(paramMap)
	}

	// Execute middleware chain
	finalHandler := handler
	for i := len(r.middleware) - 1; i >= 0; i-- {
		finalHandler = r.middleware[i](finalHandler)
	}

	// Execute handler
	finalHandler(w, req, ctx)
}

// Use adds middleware to the router
func (r *Router) Use(middleware ...Middleware[Context]) {
	r.middleware = append(r.middleware, middleware...)
}

// Group creates a new route group
func (r *Router) Group(prefix string) *Group {
	group := &Group{
		router:     r,
		prefix:     prefix,
		middleware: []Middleware[Context]{},
	}
	r.groups = append(r.groups, group)
	return group
}

// Use adds middleware to the group
func (g *Group) Use(middleware ...Middleware[Context]) *Group {
	g.middleware = append(g.middleware, middleware...)
	return g
}

// SubGroup creates a sub-group
func (g *Group) SubGroup(prefix string) *Group {
	return &Group{
		router:     g.router,
		prefix:     g.prefix + prefix,
		middleware: append([]Middleware[Context]{}, g.middleware...),
		parent:     g,
	}
}

// Group creates a new sub-group (alias for SubGroup)
func (g *Group) Group(prefix string) *Group {
	return g.SubGroup(prefix)
}

// Register adds a route with custom method
func (r *Router) Register(method, path string, handler HandlerFunc[Context]) {
	r.tree.Add(method, path, handler)
}

// Register adds a route with custom method to group
func (g *Group) Register(method, path string, handler HandlerFunc[Context]) {
	fullPath := g.prefix + path

	// Create a wrapper that applies group middleware
	wrappedHandler := func(w http.ResponseWriter, r *http.Request, ctx Context) {
		// Apply group middleware in order
		finalHandler := handler
		for i := len(g.middleware) - 1; i >= 0; i-- {
			finalHandler = g.middleware[i](finalHandler)
		}
		finalHandler(w, r, ctx)
	}

	g.router.tree.Add(method, fullPath, wrappedHandler)
}

// HTTP method helpers for Router
func (r *Router) GET(path string, handler HandlerFunc[Context]) {
	r.Register(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler HandlerFunc[Context]) {
	r.Register(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler HandlerFunc[Context]) {
	r.Register(http.MethodPut, path, handler)
}

func (r *Router) PATCH(path string, handler HandlerFunc[Context]) {
	r.Register(http.MethodPatch, path, handler)
}

func (r *Router) DELETE(path string, handler HandlerFunc[Context]) {
	r.Register(http.MethodDelete, path, handler)
}

func (r *Router) HEAD(path string, handler HandlerFunc[Context]) {
	r.Register(http.MethodHead, path, handler)
}

func (r *Router) OPTIONS(path string, handler HandlerFunc[Context]) {
	r.Register(http.MethodOptions, path, handler)
}

// DELEGATE creates a delegate route (for static file serving)
func (r *Router) DELEGATE(path string, method string, handler HandlerFunc[Context]) {
	r.Register(method, path, handler)
}

// HTTP method helpers for Group
func (g *Group) GET(path string, handler HandlerFunc[Context]) {
	g.Register(http.MethodGet, path, handler)
}

func (g *Group) POST(path string, handler HandlerFunc[Context]) {
	g.Register(http.MethodPost, path, handler)
}

func (g *Group) PUT(path string, handler HandlerFunc[Context]) {
	g.Register(http.MethodPut, path, handler)
}

func (g *Group) PATCH(path string, handler HandlerFunc[Context]) {
	g.Register(http.MethodPatch, path, handler)
}

func (g *Group) DELETE(path string, handler HandlerFunc[Context]) {
	g.Register(http.MethodDelete, path, handler)
}

func (g *Group) HEAD(path string, handler HandlerFunc[Context]) {
	g.Register(http.MethodHead, path, handler)
}

func (g *Group) OPTIONS(path string, handler HandlerFunc[Context]) {
	g.Register(http.MethodOptions, path, handler)
}

// DELEGATE creates a delegate route for group
func (g *Group) DELEGATE(path string, method string, handler HandlerFunc[Context]) {
	g.Register(method, path, handler)
}

// Static serves static files
func (r *Router) Static(prefix, root string) {
	handler := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
	r.DELEGATE(prefix+"*", http.MethodGet, func(w http.ResponseWriter, r *http.Request, ctx Context) {
		handler.ServeHTTP(w, r)
	})
}

// Static serves static files for group
func (g *Group) Static(prefix, root string) {
	handler := http.StripPrefix(g.prefix+prefix, http.FileServer(http.Dir(root)))
	g.DELEGATE(prefix+"*", http.MethodGet, func(w http.ResponseWriter, r *http.Request, ctx Context) {
		handler.ServeHTTP(w, r)
	})
}
