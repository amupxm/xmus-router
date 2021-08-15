package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerMethods(t *testing.T) {
	rt := NewRouter(&RouterOptions{})
	testTable := []struct {
		Method         string
		Handler        func(ctx *XmusContext)
		HandlerHandler func(path string, f func(ctx *XmusContext)) *route
	}{
		{"GET", func(ctx *XmusContext) { ctx.Response.Write([]byte("GET")) }, rt.GET},
		{"POST", func(ctx *XmusContext) { ctx.Response.Write([]byte("POST")) }, rt.POST},
		{"PUT", func(ctx *XmusContext) { ctx.Response.Write([]byte("PUT")) }, rt.PUT},
		{"DELETE", func(ctx *XmusContext) { ctx.Response.Write([]byte("DELETE")) }, rt.DELETE},
		{"PATCH", func(ctx *XmusContext) { ctx.Response.Write([]byte("PATCH")) }, rt.PATCH},
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
	rt := NewRouter(&RouterOptions{})
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
		rt.GET("/", func(ctx *XmusContext) {
			ctx.SetStatus(test.StatusCode).JSON(test.Json)
		})
		rt.ServeHTTP(testReq, req)
		if testReq.Body.String() != test.JsonExpected || testReq.Code != test.StatusCode || testReq.Header().Get("Content-Type") != "application/json" {
			t.Errorf("#%d: body not equal", testCase)
			continue

		}
	}
}

func BenchmarkName(b *testing.B) {
	rt := NewRouter(&RouterOptions{EchoLogs: false})
	req, _ := http.NewRequest(MethodGet, "/", nil)
	testReq := httptest.NewRecorder()
	rt.GET("/", func(ctx *XmusContext) { ctx.Response.Write([]byte("GET")) })
	for i := 0; i < b.N; i++ {
		rt.ServeHTTP(testReq, req)
	}
}
