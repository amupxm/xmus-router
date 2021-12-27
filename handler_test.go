package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerMethods(t *testing.T) {
	rt := NewRouter(&RouterOption{})
	testTable := []struct {
		Method         string
		Handler        http.Handler
		HandlerHandler func(path string, handler http.Handler)
	}{
		{"GET", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("GET")) })
		}(), rt.GET},
		{"POST", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("POST")) })
		}(), rt.POST},
		{"PUT", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PUT")) })
		}(), rt.PUT},
		{"DELETE", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("DELETE")) })
		}(), rt.DELETE},
		{"PATCH", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PATCH")) })
		}(), rt.PATCH},
	}
	for testCase, test := range testTable {
		req, _ := http.NewRequest(test.Method, "/", nil)
		testReq := httptest.NewRecorder()
		test.HandlerHandler("/", test.Handler)
		rt.ServeHTTP(testReq, req)
		if testReq.Body.String() != test.Method {
			t.Errorf("#%d: body not equal", testCase)
			continue

		}
	}
}

func TestHandlerRegister(t *testing.T) {
	rt := NewRouter(&RouterOption{})
	testTable := []struct {
		Path    string
		Method  string
		Handler http.Handler
	}{
		{"/", "GET", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("GET")) })
		}()},
		{"/", "POST", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("POST")) })
		}()},
		{"/", "PUT", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PUT")) })
		}()},
		{"/", "DELETE", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("DELETE")) })
		}()},
		{"/", "PATCH", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PATCH")) })
		}()},
		{"/hello/", "PATCH", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PATCH")) })
		}()},
		{"/param1/param2/param3/param4/", "PATCH", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PATCH")) })
		}()},
		{"/p1/", "PATCH", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PATCH")) })
		}()},
	}
	for testCase, test := range testTable {
		req, _ := http.NewRequest(test.Method, test.Path, nil)
		testReq := httptest.NewRecorder()
		rt.Register(test.Path, test.Method, test.Handler)
		rt.ServeHTTP(testReq, req)
		if testReq.Body.String() != test.Method {
			t.Errorf("#%d: body not equal", testCase)
			continue
		}
	}
}
