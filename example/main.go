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
	rt := router.NewRouter()
	h1 := NewH("hi1")
	h2 := NewH("hi2")
	h3 := NewH("hi3")

	rt.Register(h1, "/", "GET")

	rt.Register(h2, "/hi/:dd/:cc/23/", "GET")

	rt.Register(h3, "/by/", "GET")

	http.ListenAndServe(":8080", rt)
}
func (h *handlerOne) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.text))
}
