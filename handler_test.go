package router

import (
	"encoding/json"
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

func TestContextJsonAndStatus(t *testing.T) {
	rt := NewRouter(&RouterOption{})
	type jsonType struct {
		SampleTest string `json:"sample_test"`
		SampleInt  int    `json:"sample_int"`
	}
	testTable := []struct {
		Json         jsonType
		StatusCode   int
		JsonExpected string
	}{
		{jsonType{SampleTest: "sample_test6", SampleInt: 1}, 200, `{"sample_test":"sample_test6","sample_int":1}`},
		{jsonType{SampleTest: "sample_test5", SampleInt: 2}, 500, `{"sample_test":"sample_test5","sample_int":2}`},
		{jsonType{SampleTest: "sample_test4", SampleInt: 3}, 403, `{"sample_test":"sample_test4","sample_int":3}`},
		{jsonType{SampleTest: "sample_test3", SampleInt: 4}, 404, `{"sample_test":"sample_test3","sample_int":4}`},
		{jsonType{SampleTest: "sample_test2", SampleInt: 5}, 500, `{"sample_test":"sample_test2","sample_int":5}`},
		{jsonType{SampleTest: "sample_test1", SampleInt: 6}, 501, `{"sample_test":"sample_test1","sample_int":6}`},
	}
	for testCase, test := range testTable {
		req, _ := http.NewRequest(MethodGet, "/", nil)
		testReq := httptest.NewRecorder()
		rt.GET("/", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.StatusCode)
				w.Header().Set("Content-Type", "application/json")
				marshaled, _ := json.Marshal(test.Json)
				w.Write([]byte(marshaled))
			})
		}())
		rt.ServeHTTP(testReq, req)
		if testReq.Body.String() != test.JsonExpected || testReq.Code != test.StatusCode || testReq.Header().Get("Content-Type") != "application/json" {
			t.Errorf("#%d: body not equal , expected %s , got %s", testCase, test.JsonExpected, testReq.Body.String())
			continue

		}
	}
}

func TestHttpMethoddNotAllowed(t *testing.T) {
	rt := NewRouter(&RouterOption{})

	testTable := []struct {
		Method             string
		RequestPath        string
		ExpectedStatusCode int
	}{
		{"POST", "/", 405},
		{"PUT", "/", 405},
		{"PATCH", "/", 405},
	}
	for testCase, test := range testTable {
		req, _ := http.NewRequest(test.Method, test.RequestPath, nil)
		rt.GET(test.RequestPath, func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
		}())
		testReq := httptest.NewRecorder()
		rt.ServeHTTP(testReq, req)
		if testReq.Code != test.ExpectedStatusCode {
			t.Errorf("#%d: response code is not equal , got %d , expected %d", testCase, testReq.Code, test.ExpectedStatusCode)
			continue

		}
	}
}

func TestHttpNotFound(t *testing.T) {
	rt := NewRouter(&RouterOption{})

	testTable := []struct {
		Method             string
		RequestPath        string
		ExpectedStatusCode int
	}{
		{"POST", "/hello", 404},
		{"PUT", "/bye", 404},
		{"PATCH", "/amupxm", 404},
	}
	for testCase, test := range testTable {
		req, _ := http.NewRequest(test.Method, test.RequestPath, nil)
		rt.GET("/", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
		}())
		testReq := httptest.NewRecorder()
		rt.ServeHTTP(testReq, req)
		if testReq.Code != test.ExpectedStatusCode {
			t.Errorf("#%d: response code is not equal , got %d , expected %d", testCase, testReq.Code, test.ExpectedStatusCode)
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

func TestHandlerDelegate(t *testing.T) {
	testTable := []struct {
		Path          string
		PathToRequest string
		Method        string
		Handler       http.Handler
	}{

		{"/", "/ok!/", "DELETE", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("DELETE")) })
		}()},
		{"/", "/hello/request/", "PATCH", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PATCH")) })
		}()},
		{"/hello/", "/hello/", "POST", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("POST")) })
		}()},
		{"/param1/param2/param3/param4/", "/param1/param2/param3/param4/2/param2/param1/", "GET", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("GET")) })
		}()},
		{"/p1/", "/p1/p2/313123", "KICK", func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("KICK")) })
		}()},
	}
	for testCase, test := range testTable {
		rt := NewRouter(&RouterOption{})
		req, _ := http.NewRequest(test.Method, test.PathToRequest, nil)
		testReq := httptest.NewRecorder()
		rt.DELEGATE(test.Path, test.Method, test.Handler)
		rt.ServeHTTP(testReq, req)
		if testReq.Body.String() != test.Method {
			t.Errorf("#%d: body not equal,got %v , expected %v", testCase, testReq.Body.String(), test.Method)
			continue
		}
	}
}
