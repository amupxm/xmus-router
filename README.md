# XMUS-router

Simple router build on `net/http` supports custom middleWare.

If any feature is needed, please report it as a bug.
 


## usage

simply import router and use it ✌️.


```go
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

```

the `RouterContext` contains request , responseWriter and url parameters plus a few methods to make is easy to use.

