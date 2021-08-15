package router

import "fmt"

func (r *router) GET(path string, f func(ctx *XmusContext)) *route {
	return r.AddCustomMethodRoute("GET", path, f)
}

func (r *router) POST(path string, f func(ctx *XmusContext)) *route {
	return r.AddCustomMethodRoute("POST", path, f)
}
func (r *router) PUT(path string, f func(ctx *XmusContext)) *route {
	return r.AddCustomMethodRoute("PUT", path, f)
}
func (r *router) DELETE(path string, f func(ctx *XmusContext)) *route {
	return r.AddCustomMethodRoute("DELETE", path, f)
}
func (r *router) PATCH(path string, f func(ctx *XmusContext)) *route {
	return r.AddCustomMethodRoute("PATCH", path, f)
}
func (r *router) DELEGATE(path, method string, f func(ctx *XmusContext)) *route {
	return r.AddCustomMethodRoute(method, fmt.Sprintf("%s*/", path), f)
}
