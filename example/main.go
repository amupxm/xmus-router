package main

import (
	"net/http"

	router "github.com/amupxm/xmus-router"
)

type (
	handlerOne struct {
		text string
	}
	HandlerOne interface {
		ServeHTTP(http.ResponseWriter, *http.Request)
	}
)

func NewH(test string) HandlerOne {
	return &handlerOne{test}
}
func main() {
	rt := router.NewRouter(&router.RouterOption{})
	h1 := NewH("hi1")

	rt.GET("/hello/:1/", h1)

	http.ListenAndServe(":8080", rt)
}
func (h *handlerOne) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.text))
}
