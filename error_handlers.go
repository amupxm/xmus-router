package router

import "net/http"

type (
	notFound      struct{}
	notNotAllowed struct{}
)

func (nt notFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header()["Content-Type"] = []string{"application/json"}
	w.WriteHeader(http.StatusNotFound)
	w.Write(errorNotFoundMessage)
}

func (nt notNotAllowed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header()["Content-Type"] = []string{"application/json"}
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(errorMethodNotAllowedMessage)
}
