package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_URLParamsAll(t *testing.T) {
	rt := NewRouter(&RouterOptions{EchoLogs: false})

	testTable := []struct {
		Url        string
		RequestUrl string
		Method     string
		Params     map[string]string
	}{
		{"/app/:version/:id/", "/app/1.0/1", "GET", map[string]string{"version": "1.0", "id": "1"}},
		{"/app/:version/:id/", "/app/1.0/131231/", "GET", map[string]string{"version": "1.0", "id": "131231"}},
		{"/:id/", "/131231", "GET", map[string]string{"id": "131231"}},
	}
	for testCase, test := range testTable {
		req, _ := http.NewRequest(test.Method, test.RequestUrl, nil)
		testReq := httptest.NewRecorder()
		rt.AddCustomMethodRoute(test.Method, test.Url, func(context *XmusContext) {
			marshaled, err := json.Marshal(context.URLParams)
			if err != nil {
				t.Errorf("Test %d:error on marshaling url parameters %s", testCase, err)
			}
			context.Response.Write(marshaled)
		})
		rt.ServeHTTP(testReq, req)
		marshaled, err := json.Marshal(test.Params)
		if err != nil {
			t.Errorf("Test %d:error on marshaling url parameters %s", testCase, err)
		}
		if testReq.Body.String() != string(marshaled) {
			t.Errorf("#%d: body not equal, got %s , expected %s", testCase, testReq.Body.String(), string(marshaled))
			continue

		}
	}
}

func Test_URLParamsOne(t *testing.T) {
	rt := NewRouter(&RouterOptions{})

	testTable := []struct {
		Url        string
		RequestUrl string
		Method     string
		Key        string
		Value      string
	}{
		{"/app/:version/", "/app/1.0", "GET", "version", "1.0"},
		{"/app/:id/:name/", "/app/1/amupxm", "GET", "name", "amupxm"},
		{"/sam3ple/version/", "/sam3ple/version", "GET", "version", ""},
	}
	for testCase, test := range testTable {
		req, _ := http.NewRequest(test.Method, test.RequestUrl, nil)
		testReq := httptest.NewRecorder()
		rt.AddCustomMethodRoute(test.Method, test.Url, func(context *XmusContext) {
			context.Response.Write([]byte(context.GetParam(test.Key)))
		})
		rt.ServeHTTP(testReq, req)
		fmt.Println(test.Key, test.Value)
		if testReq.Body.String() != test.Value {
			t.Errorf("#%d: body not equal, got %s , expected %s", testCase, testReq.Body.String(), test.Value)
			continue
		}
	}
}
