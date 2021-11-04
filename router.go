package router

import (
	"net/http"
	"strings"
)

type (
	Router interface {
		ServeHTTP(http.ResponseWriter, *http.Request)
		Register(path, method string, handler http.Handler) error
		GET(path string, handler http.Handler)
		POST(path string, handler http.Handler)
		PUT(path string, handler http.Handler)
		DELETE(path string, handler http.Handler)
		PATCH(path string, handler http.Handler)
		DELEGATE(path string, method string, handler http.Handler)
	}
	router struct {
		notFoundHandler  http.Handler
		methodNotAllowed http.Handler
		routes           map[string]routeGroup

		logf func(format string, args ...interface{})
	}
	routeGroup struct {
		routes     []route
		pathArr    []string
		params     []string
		isDelegate bool
		hasParams  bool
	}
	route struct {
		handler http.Handler
		method  string
	}
	middleware struct {
		handler http.Handler
		next    http.Handler
	}

	RouterOption struct {
		NotFoundHandler  http.Handler
		MethodNotAllowed http.Handler
		Logf             func(format string, args ...interface{})
	}
)

func NewRouter(opts *RouterOption) Router {
	var notFoundHandler notFound
	var methodNotAllowedHandler notNotAllowed

	r := router{
		notFoundHandler:  notFoundHandler,
		methodNotAllowed: methodNotAllowedHandler,
		routes:           make(map[string]routeGroup),
	}
	if opts.MethodNotAllowed != nil {
		r.notFoundHandler = opts.MethodNotAllowed
	}
	if opts.NotFoundHandler != nil {
		r.notFoundHandler = opts.NotFoundHandler
	}
	if nil != opts.Logf {
		r.logf = opts.Logf
	}
	return &r
}

func (rt *router) Register(path, method string, handler http.Handler) error {
	// validate path
	path = validatePath(path)
	if _, ok := rt.routes[path]; !ok {
		hashParams, isDelegate, URLParams := getPathInfo(path)
		rt.routes[path] = routeGroup{
			routes:     make([]route, 0),
			pathArr:    strings.Split(path, "/"),
			params:     URLParams,
			isDelegate: isDelegate,
			hasParams:  hashParams,
		}
	}
	// create route
	route := route{
		handler: handler,
		method:  method,
	}

	// add route to router
	rTemp := rt.routes[path]
	rTemp.routes = append(rt.routes[path].routes, route)
	rt.routes[path] = rTemp

	if rt.logf != nil {
		rt.logf("Path : %s with method %s regstered", path, method)
	}
	return nil
}

func (rt router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// prepare request path
	reqPath := prepareRequestPath(r.URL.Path)
	// get routes
	routes, err := matchPath(rt.routes, reqPath)
	if err != nil {
		rt.notFoundHandler.ServeHTTP(w, r) // TODO  : logf request and responser
		return

	}
	// get handler
	route, err := rt.matchMethod(routes, r.Method)
	if err != nil {
		rt.methodNotAllowed.ServeHTTP(w, r) // TODO  : logf request and responser
		return
	}
	route.handler.ServeHTTP(w, r) // TODO  : logf request and responser
}
func (rt router) matchMethod(r []route, method string) (re route, err error) {
	for _, route := range r {
		if route.method == method {
			return route, nil
		}
	}
	return re, errMethodNotAllowed
}

func matchPath(routes map[string]routeGroup, reqPath string) ([]route, error) {
	splitedReq := strings.Split(reqPath, "/")
	for routePath, routegp := range routes {
		if reqPath == routePath {
			return routegp.routes, nil
		}
		if !routegp.isDelegate && !routegp.hasParams && reqPath != routePath {
			continue
		}
		// if has params and delegate
		if (routegp.isDelegate) || (routegp.hasParams && len(splitedReq) == len(routegp.pathArr)) {
			checkMatched := func() bool {
				for i, v := range routegp.pathArr {
					if v == "*" {
						return true
					}
					if i > len(splitedReq)-1 {
						return false
					}
					if v != splitedReq[i] && !isParamKey(routegp.params, v) {
						return false
					}
				}
				return true
			}()
			if checkMatched {
				return routegp.routes, nil
			}
		}
	}
	return nil, errNotFound
}

func isParamKey(params []string, key string) bool {
	for _, v := range params {
		if len(key) <= 1 {
			return false
		}
		if v == key[1:] {
			return true
		}
	}
	return false
}
