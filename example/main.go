package example

import (
	"net/http"

	router "github.com/amupxm/xmus-router"
)

func main() {
	router := router.NewRouter(RouterConfig{})
	group1 := router.NewGroup().addMiddleware(middlewareFunc1).addMiddleware(middlewareFunc2)
	group1.AddGroup("/v1")
	group1.Get("/api").Register(http.HandlerFunc)
	group2 := group1.AddGroup("/posts")
	group2.Post("/").Register(http.HandleFunc)
	group2.Update("/{postID}").Register(http.HandleFunc)

	http.ListenAndServe(":8080", router)
}
