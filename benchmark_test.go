package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Benchmark5Params(b *testing.B) {
	rt := NewRouter(&RouterOption{})
	req, _ := http.NewRequest(MethodGet, "/param/path/to/parameter/john/12345", nil)
	testReq := httptest.NewRecorder()
	rt.Register("/param/:param1/:params2/:param3/:param4/:param5/", "GET", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	for i := 0; i < b.N; i++ {
		rt.ServeHTTP(testReq, req)
	}
}

func BenchmarkOneRoute(b *testing.B) {
	rt := NewRouter(&RouterOption{})
	req, _ := http.NewRequest(MethodGet, "/param", nil)
	testReq := httptest.NewRecorder()
	rt.Register("/param/", "GET", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	for i := 0; i < b.N; i++ {
		rt.ServeHTTP(testReq, req)
	}
}
