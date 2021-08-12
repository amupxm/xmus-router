package main

import (
	"log"
	"net/http"

	router "github.com/amupxm/xmus-router"
	middlewareLogger "github.com/amupxm/xmus-router/middleware/logger"
)

func main() {
	router := router.NewRouter()
	buildInLogger := middlewareLogger.Logger
	router.CustomMethodRequest("GET", "/hello/:id/:user/", SampleHandlerNumOne)
	router.GET("/log/user/agent/", SampleHandlerNumTwo).AddMiddleWare(buildInLogger).AddMiddleWare(LogUserAgent)
	http.ListenAndServe(":8080", router)
}
func SampleHandlerNumOne(c *router.RouterContext) {
	c.SetStatus(http.StatusOK).SetHeader("Content-Type", "text/html").
		JSON(map[string]string{"id": c.URLParams["id"], "user": c.URLParams["user"]})
}

func SampleHandlerNumTwo(c *router.RouterContext) {
	c.SetStatus(http.StatusOK).SetHeader("Content-Type", "text/html")
	c.Response.Write([]byte("<html><h1><b>Hello World!</b></h1></html>"))
}

func LogUserAgent(c *router.RouterContext) {
	log.Println(c.Request.UserAgent())
}
