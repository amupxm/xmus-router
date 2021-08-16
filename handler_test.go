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
