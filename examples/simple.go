package main

import (
	"fmt"
	"log"
	"net/http"

	rtr "github.com/amupxm/xmus-router"
)

func main() {
	router := rtr.NewRouter(&rtr.RouterOptions{})
	router.AddCustomMethodRoute("GET", "/yes/", SampleHandlerNumOne)
	router.AddCustomMethodRoute("GET", "/yes/:id/:Lp/", SampleHandlerNumOne)
	router.GET("/log/", SampleHandlerNumOne)
	router.DELEGATE("/dd/", rtr.MethodGet, SampleHandlerNumOne)
	// router.CustomMethodRequest("GET", "/hello/:id/:user/", SampleHandlerNumOne)
	// router.GET("/log/user/agent/", SampleHandlerNumTwo).AddMiddleWare(buildInLogger).AddMiddleWare(LogUserAgent)
	http.ListenAndServe(":8080", router)
}
func SampleHandlerNumOne(c *rtr.XmusContext) {
	fmt.Println(c.GetParam("Lp"))
	c.Response.Write([]byte("yes"))
	// c.SetStatus(http.StatusOK).SetHeader("Content-Type", "text/html").
	// 	JSON(map[string]string{"id": c.URLParams["id"], "user": c.URLParams["user"]})
}

// func SampleHandlerNumTwo(c *router.RouterContext) {
// 	c.SetStatus(http.StatusOK).SetHeader("Content-Type", "text/html")
// 	c.Response.Write([]byte("<html><h1><b>Hello World!</b></h1></html>"))
// }

func LogUserAgent(c *rtr.XmusContext) {
	log.Println(c.Request.UserAgent())
}
