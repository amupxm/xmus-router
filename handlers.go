package router

func (r *router) GET(path string, f func(ctx *RouterContext)) {
	r.CustomMethodRequest("GET", path, f)
}
func (r *router) POST(path string, f func(ctx *RouterContext)) {
	r.CustomMethodRequest("POST", path, f)
}
func (r *router) PUT(path string, f func(ctx *RouterContext)) {
	r.CustomMethodRequest("PUT", path, f)
}
func (r *router) DELETE(path string, f func(ctx *RouterContext)) {
	r.CustomMethodRequest("DELETE", path, f)
}
func (r *router) PATCH(path string, f func(ctx *RouterContext)) {
	r.CustomMethodRequest("PATCH", path, f)
}
func (r *router) OPTIONS(path string, f func(ctx *RouterContext)) {
	r.CustomMethodRequest("OPTIONS", path, f)
}
