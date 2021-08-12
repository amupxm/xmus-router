package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type (
	route struct {
		Method      string
		HandlerFunc func(context *RouterContext)
		URLParams   map[string]string
		Ctx         *context.Context
		Middleware  []*Middleware
	}
	Middleware struct {
		Handler func(context *RouterContext)
	}
	RouterContext struct {
		Response  http.ResponseWriter
		Request   *http.Request
		URLParams map[string]string
	}
	router struct {
		NotFoundHandler         func(context *RouterContext)
		MethodNotAllowedHandler func(context *RouterContext)
		routes                  map[string][]*route // path : sroute map
		Handler                 func(http.ResponseWriter, *http.Request)
	}
	Router interface {
		PrepareURLParams(path string) (map[string]string, error)
		ServeHTTP(w http.ResponseWriter, r *http.Request)
		CustomMethodRequest(method, path string, f func(ctx *RouterContext)) *route
	}
)

func NewRouter() Router {
	return &router{
		routes: make(map[string][]*route),
	}
}

func (rt *router) trimPath(path string) (string, error) {
	// trim path
	path = strings.TrimSpace(path)
	// path should start with /
	if !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("path should start with /")
	}
	// path should end with /
	if !strings.HasSuffix(path, "/") {
		return "", fmt.Errorf("path should end with /")
	}
	return path, nil
}

func (rt *router) CustomMethodRequest(method, path string, f func(ctx *RouterContext)) *route {
	path, err := rt.trimPath(path)
	if err != nil {
		panic(err)
	}
	urlParams, err := rt.PrepareURLParams(path)
	if err != nil {
		panic(err)
	}
	route := rt.addRoute(path, method, urlParams, f)

	return route
}
func (r *route) AddMiddleWare(f func(context *RouterContext)) {
	r.Middleware = append(r.Middleware, &Middleware{f})
}
func (rt *router) addRoute(path, method string, urlParams map[string]string, f func(context *RouterContext)) *route {
	// check path exists then if path exists and methods are equal throw an error
	if exist := rt.routes[path]; exist != nil {
		for _, i := range exist {
			if i.Method == method {
				panic(fmt.Sprintf("duplicated route %s with method %s ", path, method))
			}
		}
	}
	route := route{Method: method, HandlerFunc: f}
	rt.routes[path] = append(rt.routes[path], &route)
	return &route
}

func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for routePath, route := range rt.routes {
		if rt.isMatchedPath(routePath, r.URL.Path) {
			for _, i := range route {
				if i.Method == r.Method {
					log.Println(2)

					urlParams, _ := rt.extractUrlParams(r.URL.Path, routePath)

					log.Println(3)

					for _, j := range i.Middleware {
						j.Handler(&RouterContext{Response: w, Request: r, URLParams: urlParams})
					}
					i.HandlerFunc(&RouterContext{Response: w, Request: r, URLParams: urlParams})
					return
				}
			}
			MethodNotAllowed(&RouterContext{Response: w, Request: r, URLParams: map[string]string{}})
			return
		}
	}
	NotFoundHandler(&RouterContext{Response: w, Request: r, URLParams: map[string]string{}})
}
