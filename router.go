package router

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type (
	Router interface {
		ServeHTTP(http.ResponseWriter, *http.Request)
		Register(path, method string, handler http.Handler)
		GET(path string, handler http.Handler)
		POST(path string, handler http.Handler)
		PUT(path string, handler http.Handler)
		DELETE(path string, handler http.Handler)
		PATCH(path string, handler http.Handler)
	}
	router struct {
		notFoundHandler  http.Handler
		methodNotAllowed http.Handler
		routes           groupOfRoutes
		routesWithParams groupOfRoutes
		logf             LeveledLoggerInterface
	}

	groupOfRoutes map[Path]map[Method]http.Handler

	Path         string
	Method       string
	RouterOption struct {
		NotFoundHandler  http.Handler
		MethodNotAllowed http.Handler
		Logf             LeveledLoggerInterface
	}
)

func NewRouter(opts *RouterOption) Router {
	var notFoundHandler notFound
	var methodNotAllowedHandler notNotAllowed

	r := router{
		notFoundHandler:  notFoundHandler,
		methodNotAllowed: methodNotAllowedHandler,
		routes:           make(groupOfRoutes),
	}
	if opts == nil || opts.MethodNotAllowed != nil {
		r.notFoundHandler = opts.MethodNotAllowed
	}
	if opts == nil || opts.NotFoundHandler != nil {
		r.notFoundHandler = opts.NotFoundHandler
	}
	// if opts == nil || nil != opts.Logf {
	// 	r.logf = opts.Logf
	// }
	r.routes = groupOfRoutes{}
	r.routesWithParams = groupOfRoutes{}
	return &r
}

var ErrRouteNotFound = errors.New("route not found")

func (rt *router) Register(p, m string, handler http.Handler) {
	path := Path(p)
	method := Method(m)
	path.Validate()
	// if its param route
	if strings.ContainsAny(path.String(), ":") {
		//register with params
		//replace every word begans with : with *
		arr := strings.Split(path.String(), "/")
		for i := 0; i < len(arr); i++ {
			if strings.HasPrefix(arr[i], ":") {
				arr[i] = "*"
			}
		}
		path = Path(strings.Join(arr, "/"))
		t := rt.routesWithParams
		if _, ok := t[Path(path)][Method(method)]; ok {
			panic(fmt.Sprintf("route %s with method %s already registered", path, method))
		}
		if t[Path(path)] == nil {
			t[Path(path)] = make(map[Method]http.Handler)
		}
		t[Path(path)][Method(method)] = handler
		rt.routesWithParams = t
	} else {
		t := rt.routes
		if _, ok := t[Path(path)][Method(method)]; ok {
			panic(fmt.Sprintf("route %s with method %s already registered", path, method))
		}
		if t[Path(path)] == nil {
			t[Path(path)] = make(map[Method]http.Handler)
		}

		t[Path(path)][Method(method)] = handler
		rt.routes = t
	}
}

func (rt router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	if reqPath == "" {
		reqPath = "/"
	}
	if reqPath != "/" && len(reqPath) > 1 {
		if !validateRequestPathRegex.MatchString(reqPath) {
			reqPath = fmt.Sprintf("%s/", reqPath)
		}
	}

	// 1 check main routes
	if handler, ok := rt.routes[Path(reqPath)][Method(r.Method)]; ok {
		handler.ServeHTTP(w, r)
		return
	}
	// 2 check routes with params
	for path, handlers := range rt.routesWithParams {
		splicedReq := strings.Split(reqPath, "/")
		splicedPath := strings.Split(path.String(), "/")
		if len(splicedReq) != len(splicedPath) {
			continue
		}
		ok := true
		for i := 0; i < len(splicedReq); i++ {
			if splicedPath[i] == "*" || splicedReq[i] == splicedPath[i] {
				continue
			} else {
				ok = false
				break
			}
		}
		if ok {
			handler := handlers[Method(r.Method)]
			if nil != handler {
				handler.ServeHTTP(w, r)
				return
			} else {
				rt.methodNotAllowed.ServeHTTP(w, r)
				return
			}
		}
	}
	rt.notFoundHandler.ServeHTTP(w, r)

}

// 	// // prepare request path
// 	// reqPath := prepareRequestPath(r.URL.Path)
// 	// // get routes
// 	// routes, err := matchPath(rt.routes, reqPath)
// 	// if err != nil {
// 	// 	rt.notFoundHandler.ServeHTTP(w, r) // TODO  : logf request and responser
// 	// 	return

// 	// }
// 	// // get handler
// 	// route, err := rt.matchMethod(routes, r.Method)
// 	// if err != nil {
// 	// 	rt.methodNotAllowed.ServeHTTP(w, r) // TODO  : logf request and responser
// 	// 	return
// 	// }
// 	// route.handler.ServeHTTP(w, r) // TODO  : logf request and responser
// }

// func (rt router) matchMethod(r []route, method string) (re route, err error) {
// 	for _, route := range r {
// 		if route.method == method {
// 			return route, nil
// 		}
// 	}
// 	return re, errMethodNotAllowed
// }

// func matchPath(routes map[string]routeGroup, reqPath string) ([]route, error) {
// 	splitedReq := strings.Split(reqPath, "/")
// 	for routePath, routegp := range routes {
// 		if reqPath == routePath {
// 			return routegp.routes, nil
// 		}
// 		if !routegp.isDelegate && !routegp.hasParams && reqPath != routePath {
// 			continue
// 		}
// 		// if has params and delegate
// 		if (routegp.isDelegate) || (routegp.hasParams && len(splitedReq) == len(routegp.pathArr)) {
// 			checkMatched := func() bool {
// 				for i, v := range routegp.pathArr {
// 					if v == "*" {
// 						return true
// 					}
// 					if i > len(splitedReq)-1 {
// 						return false
// 					}
// 					if v != splitedReq[i] && !isParamKey(routegp.params, v) {
// 						return false
// 					}
// 				}
// 				return true
// 			}()
// 			if checkMatched {
// 				return routegp.routes, nil
// 			}
// 		}
// 	}
// 	return nil, errNotFound
// }
