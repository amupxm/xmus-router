package router

import (
	"testing"
)

func TestValidatePath_Success(t *testing.T) {
	testTable := []struct {
		P, R string
	}{
		{"", "/"},
		{"/", "/"},
		{"/a/", "/a/"},
		{"/a/a/", "/a/a/"},
	}
	for testCase, test := range testTable {
		if path := validatePath(test.P); path != test.R {
			t.Errorf("#%d failed: got %s , expected %s", testCase, path, test.R)
			continue
		}
	}
}
func TestValidatePath_Failed(t *testing.T) {
	testTable := []struct {
		P string
	}{
		{"/a"},
		{"/a/a"},
		{"/a/a//"},
		{"/a/a/:a/:a/"},
	}
	for testCase, test := range testTable {
		//check any panic
		defer func() {
			if errCase := recover(); errCase == nil {
				t.Errorf("#%d : expected a panic but nothing happend ", testCase) // to prevent uninitialized panic
			}
		}()
		_ = validatePath(test.P)
	}
}

func TestPrepareRequestPath(t *testing.T) {
	testTable := []struct {
		P, R string
	}{
		{"", "/"},
		{"/", "/"},
		{"/a", "/a/"},
		{"/a/", "/a/"},
		{"/a/a", "/a/a/"},
		{"/a/a/", "/a/a/"},
	}
	for testCase, test := range testTable {
		//check any panic
		if p := prepareRequestPath(test.P); p != test.R {
			t.Errorf("#%d failed: got %s , expected %s", testCase, p, test.R)
			continue
		}
	}
}

func TestGetPathInfo(t *testing.T) {

	testTable := []struct {
		path                  string
		hasParams, isDelegate bool
		URLParams             []string
	}{
		{"/", false, false, nil},
		{"/a/", false, false, nil},
		{"/:a/", true, false, []string{"a"}},
		{"/:a/b/", true, false, []string{"a"}},
		{"/:a/:b/", true, false, []string{"a", "b"}},
		{"/:a/:b/c/", true, false, []string{"a", "b"}},
		{"/:a/:b/:c/", true, false, []string{"a", "b", "c"}},
		{"/a/:b/:c/", true, false, []string{"b", "c"}},
		{"/a/b/:c/", true, false, []string{"c"}},
		{"/a/b/:cc/", true, false, []string{"cc"}},
		{"/a/:cb/:c/", true, false, []string{"cb", "c"}},
		{"/a/b/c", false, false, nil},

		/// Exptact delegate
		{"/a/:cb/:c/", true, false, []string{"cb", "c"}},
		{"/a/:cb/*/", true, true, []string{"cb"}},
		{"/a/*/:c/", true, false, []string{"c"}},
		{"/a/asd/*/", false, true, nil},
	}
	for testCase, test := range testTable {
		hasParams, isDelegate, URLParams := getPathInfo(test.path)
		if hasParams != test.hasParams || isDelegate != test.isDelegate {
			t.Errorf("#%d failed: got %v, %v, %v , expected %v, %v, %v", testCase, hasParams, isDelegate, URLParams, test.hasParams, test.isDelegate, test.URLParams)
			continue
		}
		ln := 0
		for _, v1 := range URLParams {
			for _, v2 := range test.URLParams {
				if v1 == v2 {
					ln++
				}
			}
		}
		if ln != len(test.URLParams) {
			t.Errorf("#%d failed: got %v, %v, %v , expected %v, %v, %v", testCase, hasParams, isDelegate, URLParams, test.hasParams, test.isDelegate, test.URLParams)
			continue

		}
	}
}
