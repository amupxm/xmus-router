package router

import (
	"fmt"
	"net/http"
)

func (rt router) GET(path string, handler http.Handler) {
	rt.Register(path, http.MethodGet, handler)
}
func (rt router) POST(path string, handler http.Handler) {
	rt.Register(path, http.MethodPost, handler)
}
func (rt router) PUT(path string, handler http.Handler) {
	rt.Register(path, http.MethodPut, handler)
}
func (rt router) DELETE(path string, handler http.Handler) {
	rt.Register(path, http.MethodDelete, handler)
}
func (rt router) PATCH(path string, handler http.Handler) {
	rt.Register(path, http.MethodPatch, handler)
}
func (rt router) DELEGATE(path string, method string, handler http.Handler) {
	rt.Register(fmt.Sprintf("%s*/", path), method, handler)
}
