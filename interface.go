package router

import "net/http"

// Type-safe context
type Context interface {
	Request() *http.Request
	Response() ResponseWriter
	Param(key string) string
	Query(key string) string
	Set(key string, value any)
	Get(key string) (any, bool)
}

// Generic handlers
type Handler[T Context] interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, ctx T)
}

type Middleware[T Context] func(Handler[T]) Handler[T]

// ResponseWriter interface for enhanced response handling
type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	Size() int
	Written() bool
}

// Context implementation

type xmusContext struct {
	request  *http.Request
	response ResponseWriter
	params   map[string]string
	query    map[string]string
	values   map[string]any
}

func NewContext(r *http.Request) Context {
	return &xmusContext{}
}

// Implement Context interface
func (c *xmusContext) Request() *http.Request {
	return c.request
}

func (c *xmusContext) Response() ResponseWriter {
	return c.response
}

func (c *xmusContext) Param(key string) string {
	// TODO: Implement later
	return c.Param(key)
}

func (c *xmusContext) Query(key string) string {
	// TODO: Implement later
	return c.Query(key)
}

func (c *xmusContext) Set(key string, value any) {
	// TODO : implement later
	c.Set(key, value)
}

func (c *xmusContext) Get(key string) (any, bool) {
	// TODO: Implement later
	return c.Get(key)
}
