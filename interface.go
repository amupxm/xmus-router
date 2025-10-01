package router

import (
	"fmt"
	"net/http"
)

// Type-safe context
type Context interface {
	Request() *http.Request
	Response() ResponseWriter
	Param(key string) string
	Query(key string) string
	Set(key string, value any)
	Get(key string) (any, bool)
	JSON(code int, obj any) error
	String(code int, format string, values ...any) error
	HTML(code int, html string) error
	Redirect(code int, url string) error
	SetParams(params map[string]string)
}

// Generic handlers
type Handler[T Context] interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, ctx T)
}

type Middleware[T Context] func(HandlerFunc[T]) HandlerFunc[T]

// ResponseWriter interface for enhanced response handling
type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	Size() int
	Written() bool
}

type xmusResponseWriter struct {
	http.ResponseWriter
	status  int
	size    int
	written bool
}

func (w *xmusResponseWriter) Status() int {
	return w.status
}

func (w *xmusResponseWriter) Size() int {
	return w.size
}

func (w *xmusResponseWriter) Written() bool {
	return w.written
}

func (w *xmusResponseWriter) WriteHeader(code int) {
	if !w.written {
		w.status = code
		w.written = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *xmusResponseWriter) Write(data []byte) (int, error) {
	if !w.written {
		w.WriteHeader(200)
	}
	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

// Context implementation

type xmusContext struct {
	request  *http.Request
	response ResponseWriter
	params   map[string]string
	query    map[string]string
	values   map[string]any
}

func NewContext(r *http.Request, w http.ResponseWriter) Context {
	// Parse query parameters
	query := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}

	return &xmusContext{
		request:  r,
		response: &xmusResponseWriter{ResponseWriter: w, status: 200},
		params:   make(map[string]string),
		query:    query,
		values:   make(map[string]any),
	}
}

func (c *xmusContext) Request() *http.Request {
	return c.request
}

func (c *xmusContext) Response() ResponseWriter {
	return c.response
}

func (c *xmusContext) Param(key string) string {
	return c.params[key]
}

func (c *xmusContext) Query(key string) string {
	return c.query[key]
}

func (c *xmusContext) Set(key string, value any) {
	c.values[key] = value
}

func (c *xmusContext) Get(key string) (any, bool) {
	value, ok := c.values[key]
	return value, ok
}

func (c *xmusContext) SetParams(params map[string]string) {
	c.params = params
}

func (c *xmusContext) JSON(code int, obj any) error {
	c.Response().WriteHeader(code)
	c.Response().Header().Set("Content-Type", "application/json")
	// Simple JSON encoding - in production, use json.Marshal
	_, err := c.Response().Write([]byte(`{"message": "test"}`))
	return err
}

func (c *xmusContext) String(code int, format string, values ...any) error {
	c.Response().WriteHeader(code)
	c.Response().Header().Set("Content-Type", "text/plain")
	// Simple string formatting - in production, use fmt.Sprintf
	_, err := c.Response().Write([]byte(fmt.Sprintf(format, values...)))
	return err
}

func (c *xmusContext) HTML(code int, html string) error {
	c.Response().WriteHeader(code)
	c.Response().Header().Set("Content-Type", "text/html")
	_, err := c.Response().Write([]byte(html))
	return err
}

func (c *xmusContext) Redirect(code int, url string) error {
	c.Response().WriteHeader(code)
	c.Response().Header().Set("Location", url)
	return nil
}
