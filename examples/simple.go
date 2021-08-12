package main

import (
	"net/http"

	router "github.com/amupxm/xmus-router"
	connectionUtils "github.com/amupxm/xmus-router/middleware/connectionutils"
	middlewareLogger "github.com/amupxm/xmus-router/middleware/logger"
)

func main() {
	router := router.NewRouter()
	m := middlewareLogger.Logger
	m2 := connectionUtils.CancelAll
	router.CustomMethodRequest("GET", "/hello/:id/:user/", cc).AddMiddleWare(m).AddMiddleWare(m2)

	http.ListenAndServe(":8080", router)
}

func cc(r *router.RouterContext) {
	r.Response.Write([]byte("Hello World!"))
}
