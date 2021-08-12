package main

import (
	"net/http"

	router "github.com/amupxm/xmus-router"
	middlewareLogger "github.com/amupxm/xmus-router/middleware/logger"
)

func main() {
	router := router.NewRouter()
	m := middlewareLogger.Logger
	router.CustomMethodRequest("GET", "/hello/:id/:user/", cc).AddMiddleWare(m)

	http.ListenAndServe(":8080", router)
}

func cc(r *router.RouterContext) {
	r.Response.Write([]byte("Hello World!"))
}
