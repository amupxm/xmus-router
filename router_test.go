package router

// func TestMatchPath(t *testing.T) {
// 	testTable := []struct {
// 		Path          string          // path which will be register to the router
// 		RequestPathes map[string]bool // {path: willBeAccept?}
// 	}{
// 		{"/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": false, "/p1/": false, "/1/": true}},
// 		{"/p1/", map[string]bool{"/p11/": true, "/p1/p2/": false, "/p1/p2/p3/": false, "/": false}},
// 		{"/p1/p2/", map[string]bool{"/p1/p2/2/": true, "/p1/": false, "/p1/p2/p3/": false, "/": false}},
// 		{"/p1/p2/p3/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/2/": true, "/p1/p2/p3/p4/": false, "/p1/": false, "/": false}},
// 		//with oarams
// 		{"/p1/:p2/", map[string]bool{"/p1/p2/": true, "/p1/p2/p3/": false, "/p1/": false, "/": false}},
// 		{"/:p1/:p2/", map[string]bool{"/p1/p2/": true, "/p1/p2/p3/": false, "/p1/": false, "/": false}},
// 		{"/:p1/:p2/:p3/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": true, "/p1/p2/p3/p4/": false, "/p1/": false, "/": false}},
// 		{"/p1/:p2/:p3/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": true, "/p1/p2/p3/p4/": false, "/p1/": false, "/": false}},
// 		{"/p1/p2/:p3/", map[string]bool{"/p1/p2/": false, "/p1/p2/p3/": true, "/p1/p2/p3/p4/": false, "/p1/": false, "/": false}},
// 	}
// 	for _, test := range testTable {

// 		router := NewRouter(&RouterOption{
// 			NotFoundHandler:  notFound{},
// 			MethodNotAllowed: notNotAllowed{},
// 		})
// 		router.Register(test.Path, http.MethodGet, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.Write([]byte(test.Path))
// 		}))

// 		for path, _ := range test.RequestPathes {
// 			router.Register(path, http.MethodGet, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				w.Write([]byte(path))
// 			}))
// 		}

// 		req := httptest.NewRequest(http.MethodGet, test.Path, nil)
// 		w := httptest.NewRecorder()
// 		router.ServeHTTP(w, req)
// 		res := w.Result()
// 		defer res.Body.Close()
// 		data, err := ioutil.ReadAll(res.Body)
// 		assert.Nil(t, err)
// 		assert.Equal(t, test.Path, string(data))
// 	}
// }
