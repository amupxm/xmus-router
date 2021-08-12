package router_test

import (
	"testing"

	router "github.com/amupxm/xmus-router"
)

func TestPrepareURLParamsSuccessCase(t *testing.T) {

	testTable := []struct {
		Path     string
		Expected map[string]string
	}{
		{"/p1/p1/p1/", map[string]string{}},
		{"/p1/:p1/p1/", map[string]string{"p1": ""}},
		{"/p1/:p1/:p2/", map[string]string{"p1": "", "p2": ""}},
	}
	testRt := router.NewRouter()

	for testIndex, testCase := range testTable {
		r, err := testRt.PrepareURLParams(testCase.Path)
		if err != nil {
			t.Errorf("Test %d: Expected no error, got %s", testIndex, err)
			continue
		}
		if len(r) != len(testCase.Expected) {
			t.Errorf("Test %d: Expected %s parameters, got %s", testIndex, testCase.Expected, r)
			continue
		}
		for k, v := range testCase.Expected {
			if r[k] != v {
				t.Errorf("Test %d: Expected %s parameter, got %s", testIndex, testCase.Expected, r)
			}
		}

	}
}

func TestPrepareURLParamsFailCase(t *testing.T) {

	testTable := []struct {
		Path string
		Ok   bool
	}{
		{"/p1/p1/p1/", true},
		{"/p1/:p1/:p1/", false},
	}
	testRt := router.NewRouter()

	for testIndex, testCase := range testTable {
		_, err := testRt.PrepareURLParams(testCase.Path)
		if err == nil && !testCase.Ok {
			t.Errorf("Test %d: Expected error, got none", testIndex)
		}
	}
}
