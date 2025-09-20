package router

import (
	"context"
	"net/http"
)

// Router wraps the radix tree with additional functionality
type Router[T Context] struct {
	tree       *radixTree[T]
	middleware []Middleware[T]
	notFound   HandlerFunc[T]
}

// NewRouter creates a new high-performance router
func NewRouter[T Context]() *Router[T] {
	return &Router[T]{
		tree:       NewRadixTree[T](),
		middleware: []Middleware[T]{},
		notFound:   nil,
	}
}

type HandlerFunc[T Context] func(w http.ResponseWriter, r *http.Request, ctx T)

func (r *Router[T]) HandlerFunc(handler http.Handler) HandlerFunc[T] {
	return HandlerFunc[T](func(w http.ResponseWriter, r *http.Request, ctx T) {
		handler.ServeHTTP(w, r)
	})
}

// GET registers a GET route
func (r *Router[T]) GET(path string, handler http.Handler) {
	handlerFunc := r.HandlerFunc(handler)
	r.tree.Add("GET", path, handlerFunc)
}

// POST registers a POST route
func (r *Router[T]) POST(path string, handler http.Handler) {
	handlerFunc := r.HandlerFunc(handler)
	r.tree.Add("POST", path, handlerFunc)
}

// PUT registers a PUT route
func (r *Router[T]) PUT(path string, handler http.Handler) {
	handlerFunc := r.HandlerFunc(handler)
	r.tree.Add("PUT", path, handlerFunc)
}

// DELETE registers a DELETE route
func (r *Router[T]) DELETE(path string, handler http.Handler) {
	handlerFunc := r.HandlerFunc(handler)
	r.tree.Add("DELETE", path, handlerFunc)
}

// PATCH registers a PATCH route
func (r *Router[T]) PATCH(path string, handler http.Handler) {
	handlerFunc := r.HandlerFunc(handler)
	r.tree.Add("PATCH", path, handlerFunc)
}

// Add registers a route with any HTTP method
func (r *Router[T]) Add(method, path string, handler http.Handler) {
	handlerFunc := r.HandlerFunc(handler)
	r.tree.Add(method, path, handlerFunc)
}

// SetNotFound sets the not found handler
func (r *Router[T]) SetNotFound(handler http.Handler) {
	handlerFunc := r.HandlerFunc(handler)
	r.notFound = handlerFunc
}

// ServeHTTP implements http.Handler interface
func (r *Router[T]) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, params := r.tree.Find(req.Method, req.URL.Path)

	if handler != nil {
		// Create context with parameters
		ctx := r.createContext(req.Context(), params)
		handler(w, req, ctx)
	} else if r.notFound != nil {
		ctx := r.createContext(req.Context(), nil)
		r.notFound(w, req, ctx)
	} else {
		http.NotFound(w, req)
	}
}

// createContext creates a new context with parameters
// This is a placeholder - implement based on your Context type
func (r *Router[T]) createContext(base context.Context, params Parameters) T {
	// This will need to be implemented based on your specific Context type
	// For now, we'll use a type assertion that will work with context.Context
	var zero T
	return zero
}
