package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type (
	middleware struct {
		handler func(context *XmusContext) bool
	}

	router struct {
		notFoundHandler         func(context *XmusContext)
		methodNotAllowedHandler func(context *XmusContext)
		routes                  map[string][]*route // path : sroute map
		//
		echoLogs bool
	}
	Router interface {
		AddCustomMethodRoute(method string, path string, f func(context *XmusContext)) *route
		ServeHTTP(w http.ResponseWriter, r *http.Request)

		DELEGATE(path, method string, f func(context *XmusContext)) *route
		GET(path string, f func(ctx *XmusContext)) *route
		POST(path string, f func(ctx *XmusContext)) *route
		PUT(path string, f func(ctx *XmusContext)) *route
		DELETE(path string, f func(ctx *XmusContext)) *route
		PATCH(path string, f func(ctx *XmusContext)) *route
	}
	RouterOptions struct {
		NotFoundHandler         func(context *XmusContext)
		MethodNotAllowedHandler func(context *XmusContext)
		EchoLogs                bool
	}
)

func NewRouter(opt *RouterOptions) Router {

	router := router{
		routes:                  make(map[string][]*route),
		notFoundHandler:         NotFoundHandler,
		methodNotAllowedHandler: methodNotAllowed,
		echoLogs:                true,
	}
	if opt.NotFoundHandler != nil {
		router.notFoundHandler = opt.NotFoundHandler
	}
	if opt.MethodNotAllowedHandler != nil {
		router.methodNotAllowedHandler = opt.MethodNotAllowedHandler
	}
	router.echoLogs = opt.EchoLogs

	return &router
}

func (r *router) AddCustomMethodRoute(method string, path string, f func(context *XmusContext)) *route {
	if err := validatePath(path); err != nil {
		panic(err)
	}
	nr := &route{method: method, handlerFunc: f, middleware: make([]*middleware, 0)} //NewRoute
	if _, ok := r.routes[path]; !ok {
		r.routes[path] = make([]*route, 0)
	}
	r.routes[path] = append(r.routes[path], nr)
	return nr
}
func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// we have three types of pathes :
	// pathes like /amupxm/some/stuff
	// pathes like /amupxm/:param1/:param2
	// pathes like /amupxm/:param1/*
	// ist better to check first its delegate path ?
	if !strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path = fmt.Sprintf("%s/", r.URL.Path)
	}
	for routePath, routeFunc := range rt.routes {
		isDelegate := delegateRegexp.MatchString(routePath)
		hasParams := hasParamsRegexp.MatchString(routePath)
		reqArr := strings.Split(r.URL.Path, "/")
		pathArr := strings.Split(routePath, "/")
		// route without any params or delegate
		if !hasParams && !isDelegate && r.URL.Path == routePath {
			rt.mathMethod(w, r, reqArr, pathArr, hasParams, isDelegate, routeFunc)
			return
		}
		hasSameLen := len(reqArr) == len(pathArr)
		// wrong route to request
		if !hasSameLen && !isDelegate {
			continue
		}
		if (hasSameLen && hasParams) || isDelegate {
			for i, path := range pathArr {
				if path != reqArr[i] {
					// if be delegte path
					if path == "*" && isDelegate {

						rt.mathMethod(w, r, reqArr, pathArr, hasParams, isDelegate, routeFunc)
						return
					}
					if path != reqArr[i] {
						break
					}
				}
			}
			rt.mathMethod(w, r, reqArr, pathArr, hasParams, isDelegate, routeFunc)
			return
		}
	}
	rt.notFoundHandler(rt.createContex(w, r))
}

func (rt router) mathMethod(w http.ResponseWriter, r *http.Request, reqArr, pathArr []string, hasParams, hasDelegate bool, routes []*route) {
	context := rt.createContex(w, r)
	for _, route := range routes {
		if r.Method == route.method {
			if hasParams {
				context.buildParams(reqArr, pathArr)
			}
			rt.runHandler(w, r, route, context)
			return

		}
	}
	rt.methodNotAllowedHandler(context)
}

func (rt router) createContex(w http.ResponseWriter, r *http.Request) *XmusContext {
	return &XmusContext{
		Response:  w,
		Request:   r,
		URLParams: make(map[string]string),
	}
}

func (rt router) runHandler(w http.ResponseWriter, r *http.Request, route *route, context *XmusContext) {
	// TODO prevent panic
	if rt.echoLogs {
		var b string
		switch r.Method {
		case "GET":
			b = greenBg
		case "POST":
			b = blueBg
		case "PUT":
			b = yellowBg
		case "DELETE":
			b = redBg
		case "PATCH":
			b = yellowBg
		default:
			b = whiteBg
		}
		log.Printf("%v%s%v | %s | %s", b, r.Method, reset, r.URL.Path, r.Proto)
	}
	// run middlewares
	ln := len(route.middleware)
	cancel := false
	for i := 0; i < ln; i++ {
		cancel = route.middleware[i].handler(context)
		if cancel {
			break
		}
	}
	if !cancel {
		route.handlerFunc(context)
	}
}
