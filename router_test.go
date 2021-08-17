package router

import (
	"strings"
	"testing"
)

func TestMatchPath(t *testing.T) {
	testTable := []struct {
		Path          string          // path which will be register to the router
		RequestPathes map[string]bool // {path: willBeAccept?}
		isDelegate    bool
		hasParams     bool
		params        []string
	}{
		{"/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": false, "/p1/": false, "/": true}, false, false, []string{}},
		{"/p1/", map[string]bool{"/p1/": true, "/p1/p2/": false, "/p1/p2/p3/": false, "/": false}, false, false, []string{}},
		{"/p1/p2/", map[string]bool{"/p1/p2/": true, "/p1/": false, "/p1/p2/p3/": false, "/": false}, false, false, []string{}},
		{"/p1/p2/p3/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": true, "/p1/p2/p3/p4/": false, "/p1/": false, "/": false}, false, false, []string{}},
		//with oarams
		{"/p1/:p2/", map[string]bool{"/p1/p2/": true, "/p1/p2/p3/": false, "/p1/": false, "/": false}, false, true, []string{"p2"}},
		{"/:p1/:p2/", map[string]bool{"/p1/p2/": true, "/p1/p2/p3/": false, "/p1/": false, "/": false}, false, true, []string{"p1", "p2"}},
		{"/:p1/:p2/:p3/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": true, "/p1/p2/p3/p4/": false, "/p1/": false, "/": false}, false, true, []string{"p1", "p2", "p3"}},
		{"/p1/:p2/:p3/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": true, "/p1/p2/p3/p4/": false, "/p1/": false, "/": false}, false, true, []string{"p2", "p3"}},
		{"/p1/p2/:p3/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": true, "/p1/p2/p3/p4/": false, "/p1/": false, "/": false}, false, true, []string{"p3"}},
		// with delegates
		{"/*/", map[string]bool{"/p1/p2/": true, "/p1/p2/p3/": true, "/p1/": true, "/": true}, true, false, []string{}},
		{"/p1/*/", map[string]bool{"/p1/p2/": true, "/p1/p2/p3/": true, "/p1/": true, "/": false}, true, false, []string{}},
		{"/p1/*/", map[string]bool{"/p1/p2/": true, "/p1/p2/p3/": true, "/p1/": false, "/": false}, true, false, []string{}},
		{"/p1/p2/*/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": true, "/p1/": false, "/": false}, true, false, []string{}},
		{"/p1/p2/*/", map[string]bool{"/p1/p2/p3/": true, "/p1/p2/p3/amir.jpeg": true, "/p1/p2/p3/amir": true, "/p1/p2/p3/amir.mp4": true, "/p1/p2/p3/@username": true}, true, false, []string{}},
		{"/p1/p2/*/", map[string]bool{"/p1/p2/p3//": true, "/p1/p2/p3/amir.jpeg.gif": true, "/p1/p2/p3/amir/user": true, "/p1/p2/p3/*": true, "/p1/p2/p3/!@#$%^&*": true}, true, false, []string{}},
		// with delegate and params
		{"/:p1/*/", map[string]bool{"/p1/p2/": true, "/p1/p2/p3/": true, "/p1/": true, "/": false}, true, true, []string{"p1"}},
		{"/:p1/*/", map[string]bool{"/p1/p2/amypxm": true, "/p1/p2/p3/amir.png": true, "/p1/path1/path2/": true, "/": false}, true, true, []string{"p1"}},
		{"/p1/:p2/*/", map[string]bool{"/p1/p2/": true, "/p1/p2/p3/": true, "/p1/": false, "/": false}, true, true, []string{"p2"}},
		{"/p1/:p2/:p3/*/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": true, "/p1/": false, "/": false, "/p1/p2/p3/amupxm/": true}, true, true, []string{"p2", "p3"}},
	}
	for testCase, test := range testTable {
		for path, accept := range test.RequestPathes {
			fakeRoute := route{handler: nil, method: "GET"}
			arr, err := matchPath(map[string]routeGroup{
				test.Path: {pathArr: strings.Split(test.Path, "/"), routes: []route{fakeRoute}, isDelegate: test.isDelegate, params: test.params, hasParams: test.hasParams},
			}, path)
			if err != nil && accept {
				t.Errorf("#%d: got an error , expected nil in path %v with request %v ! (err : %v)", testCase, test.Path, path, err)
				continue
			}
			if accept && len(arr) == 0 {
				t.Errorf("#%d: expected to %s matches with %s , but it didn't!", testCase, path, test.Path)
				continue
			}
		}
	}
}
