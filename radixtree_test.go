package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test helper functions
func multiParamHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	userID := ctx.Param("id")
	postID := ctx.Param("postId")
	ctx.String(200, "user: %s, post: %s", userID, postID)
}

// Test NewRadixTree
func TestNewRadixTree(t *testing.T) {
	tree := NewRadixTree[Context]()
	if tree == nil {
		t.Fatal("NewRadixTree returned nil")
	}
	if tree.root == nil {
		t.Fatal("Root node is nil")
	}
	if tree.root.methods == nil {
		t.Fatal("Root methods map is nil")
	}
	if tree.root.nType != static {
		t.Fatal("Root node type should be static")
	}
}

// Test Parameters methods
func TestParameters(t *testing.T) {
	params := Parameters{
		{Key: "id", Value: "123"},
		{Key: "name", Value: "test"},
	}

	// Test Get method
	value, ok := params.Get("id")
	if !ok || value != "123" {
		t.Errorf("Get('id') = %s, %v; want '123', true", value, ok)
	}

	value, ok = params.Get("name")
	if !ok || value != "test" {
		t.Errorf("Get('name') = %s, %v; want 'test', true", value, ok)
	}

	value, ok = params.Get("nonexistent")
	if ok {
		t.Errorf("Get('nonexistent') = %s, %v; want '', false", value, ok)
	}

	// Test MustGet method
	value = params.MustGet("id")
	if value != "123" {
		t.Errorf("MustGet('id') = %s; want '123'", value)
	}

	// Test MustGet panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet should panic for nonexistent key")
		}
	}()
	params.MustGet("nonexistent")
}

// Test Add method with invalid paths
func TestAddInvalidPaths(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Test empty path
	defer func() {
		if r := recover(); r == nil {
			t.Error("Add should panic for empty path")
		}
	}()
	tree.Add("GET", "", testHandler("test"))

	// Test path not starting with /
	defer func() {
		if r := recover(); r == nil {
			t.Error("Add should panic for path not starting with /")
		}
	}()
	tree.Add("GET", "invalid", testHandler("test"))
}

// Test Add method with path too long
func TestAddPathTooLong(t *testing.T) {
	tree := NewRadixTree[Context]()
	longPath := "/" + strings.Repeat("a", 1001)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Add should panic for path too long")
		}
	}()
	tree.Add("GET", longPath, testHandler("test"))
}

// Test Find method with invalid paths
func TestFindInvalidPaths(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Test empty path
	handler, params := tree.Find("GET", "")
	if handler != nil || params != nil {
		t.Error("Find should return nil for empty path")
	}

	// Test path not starting with /
	handler, params = tree.Find("GET", "invalid")
	if handler != nil || params != nil {
		t.Error("Find should return nil for path not starting with /")
	}
}

// Test parameter routes with edge cases
func TestParameterRoutesEdgeCases(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Add parameter routes with edge cases
	tree.Add("GET", "/users/:id", paramHandler)
	tree.Add("GET", "/posts/:id/comments/:commentId", multiParamHandler)

	tests := []struct {
		path     string
		want     string
		wantCode int
	}{
		{"/users/123", "param: 123", 200},
		{"/posts/456/comments/789", "user: 456, post: 789", 200},
		{"/users", "", 0},     // No match
		{"/posts/456", "", 0}, // No match
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, params := tree.Find("GET", tt.path)
			if tt.wantCode == 0 {
				if handler != nil {
					t.Errorf("Find() = %v; want nil", handler)
				}
				return
			}
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			ctx.SetParams(make(map[string]string))
			for _, p := range params {
				ctx.SetParams(map[string]string{p.Key: p.Value})
			}
			handler(w, req, ctx)

			if w.Code != tt.wantCode {
				t.Errorf("Status = %d; want %d", w.Code, tt.wantCode)
			}
			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test wildcard routes with edge cases
func TestWildcardRoutesEdgeCases(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Add wildcard routes with edge cases
	tree.Add("GET", "/static/*path", wildcardHandler)
	tree.Add("GET", "/files/*path", wildcardHandler)

	tests := []struct {
		path     string
		want     string
		wantCode int
	}{
		{"/static/css/style.css", "wildcard: css/style.css", 200},
		{"/static/js/app.js", "wildcard: js/app.js", 200},
		{"/files/documents/report.pdf", "wildcard: documents/report.pdf", 200},
		{"/static", "", 0}, // No match
		{"/files", "", 0},  // No match
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, params := tree.Find("GET", tt.path)
			if tt.wantCode == 0 {
				if handler != nil {
					t.Errorf("Find() = %v; want nil", handler)
				}
				return
			}
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			ctx.SetParams(make(map[string]string))
			for _, p := range params {
				ctx.SetParams(map[string]string{p.Key: p.Value})
			}
			handler(w, req, ctx)

			if w.Code != tt.wantCode {
				t.Errorf("Status = %d; want %d", w.Code, tt.wantCode)
			}
			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test node splitting
func TestNodeSplitting(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Add routes that will cause node splitting
	tree.Add("GET", "/test", testHandler("test"))
	tree.Add("GET", "/testing", testHandler("testing"))
	tree.Add("GET", "/tested", testHandler("tested"))

	tests := []struct {
		path string
		want string
	}{
		{"/test", "test"},
		{"/testing", "testing"},
		{"/tested", "tested"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, _ := tree.Find("GET", tt.path)
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			handler(w, req, ctx)

			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test priority ordering
func TestPriorityOrdering(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Add routes with different priorities
	tree.Add("GET", "/api", testHandler("api"))
	tree.Add("GET", "/api/users", testHandler("users"))
	tree.Add("GET", "/api/posts", testHandler("posts"))

	// Access routes to increase priority
	tree.Find("GET", "/api/posts")
	tree.Find("GET", "/api/posts")
	tree.Find("GET", "/api/users")

	// The tree should reorder based on priority
	// This is tested implicitly by ensuring the routes still work
	tests := []struct {
		path string
		want string
	}{
		{"/api", "api"},
		{"/api/users", "users"},
		{"/api/posts", "posts"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, _ := tree.Find("GET", tt.path)
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			handler(w, req, ctx)

			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test concurrent access
func TestConcurrentAccess(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Add some routes
	tree.Add("GET", "/api", testHandler("api"))
	tree.Add("GET", "/users/:id", paramHandler)

	// Test concurrent reads
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			handler, _ := tree.Find("GET", "/api")
			if handler == nil {
				t.Error("Concurrent read failed")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	// Test min function
	if min(1, 2) != 1 {
		t.Error("min(1, 2) should return 1")
	}
	if min(2, 1) != 1 {
		t.Error("min(2, 1) should return 1")
	}
	if min(1, 1) != 1 {
		t.Error("min(1, 1) should return 1")
	}

	// Test unsafeString function
	b := []byte("test")
	s := unsafeString(b)
	if s != "test" {
		t.Errorf("unsafeString() = %s; want 'test'", s)
	}

	// Test unsafeString with empty slice
	b = []byte{}
	s = unsafeString(b)
	if s != "" {
		t.Errorf("unsafeString([]byte{}) = %s; want ''", s)
	}

	// Test unsafeBytes function
	s = "test"
	b = unsafeBytes(s)
	if string(b) != "test" {
		t.Errorf("unsafeBytes() = %s; want 'test'", string(b))
	}

	// Test unsafeBytes with empty string
	s = ""
	b = unsafeBytes(s)
	if b != nil {
		t.Errorf("unsafeBytes('') = %v; want nil", b)
	}
}

// Test node type constants
func TestNodeTypes(t *testing.T) {
	if static != 0 {
		t.Error("static should be 0")
	}
	if param != 1 {
		t.Error("param should be 1")
	}
	if wildcard != 2 {
		t.Error("wildcard should be 2")
	}
}

// Test empty tree
func TestEmptyTree(t *testing.T) {
	tree := NewRadixTree[Context]()

	handler, params := tree.Find("GET", "/any")
	if handler != nil {
		t.Error("Empty tree should return nil handler")
	}
	// Note: params will be an empty slice, not nil, due to pre-allocation
	if len(params) != 0 {
		t.Error("Empty tree should return empty params")
	}
}

// Test single character routes
func TestSingleCharacterRoutes(t *testing.T) {
	tree := NewRadixTree[Context]()

	tree.Add("GET", "/a", testHandler("a"))
	tree.Add("GET", "/b", testHandler("b"))
	tree.Add("GET", "/c", testHandler("c"))

	tests := []struct {
		path string
		want string
	}{
		{"/a", "a"},
		{"/b", "b"},
		{"/c", "c"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, _ := tree.Find("GET", tt.path)
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			handler(w, req, ctx)

			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test very long paths
func TestVeryLongPaths(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Create a very long path (but under the limit)
	longPath := "/" + strings.Repeat("a", 999)
	tree.Add("GET", longPath, testHandler("long"))

	handler, _ := tree.Find("GET", longPath)
	if handler == nil {
		t.Error("Find() = nil; want handler")
		return
	}

	// Test the handler
	req := httptest.NewRequest("GET", longPath, nil)
	w := httptest.NewRecorder()
	ctx := NewContext(req, w)
	handler(w, req, ctx)

	if !strings.Contains(w.Body.String(), "long") {
		t.Errorf("Response = %s; want to contain 'long'", w.Body.String())
	}
}

// Test parameter with special characters
func TestParameterSpecialCharacters(t *testing.T) {
	tree := NewRadixTree[Context]()

	tree.Add("GET", "/users/:id", paramHandler)

	tests := []struct {
		path string
		want string
	}{
		{"/users/123", "param: 123"},
		{"/users/user-123", "param: user-123"},
		{"/users/user_123", "param: user_123"},
		{"/users/user.123", "param: user.123"},
		{"/users/user@123", "param: user@123"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, params := tree.Find("GET", tt.path)
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			ctx.SetParams(make(map[string]string))
			for _, p := range params {
				ctx.SetParams(map[string]string{p.Key: p.Value})
			}
			handler(w, req, ctx)

			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test wildcard with special characters
func TestWildcardSpecialCharacters(t *testing.T) {
	tree := NewRadixTree[Context]()

	tree.Add("GET", "/static/*path", wildcardHandler)

	tests := []struct {
		path string
		want string
	}{
		{"/static/css/style.css", "wildcard: css/style.css"},
		{"/static/js/app.min.js", "wildcard: js/app.min.js"},
		{"/static/images/logo@2x.png", "wildcard: images/logo@2x.png"},
		{"/static/files/document-2023.pdf", "wildcard: files/document-2023.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, params := tree.Find("GET", tt.path)
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			ctx.SetParams(make(map[string]string))
			for _, p := range params {
				ctx.SetParams(map[string]string{p.Key: p.Value})
			}
			handler(w, req, ctx)

			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test multiple methods on same path
func TestMultipleMethods(t *testing.T) {
	tree := NewRadixTree[Context]()

	tree.Add("GET", "/users/:id", paramHandler)
	tree.Add("POST", "/users/:id", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		ctx.String(200, "POST user: %s", ctx.Param("id"))
	})
	tree.Add("PUT", "/users/:id", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		ctx.String(200, "PUT user: %s", ctx.Param("id"))
	})
	tree.Add("DELETE", "/users/:id", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		ctx.String(200, "DELETE user: %s", ctx.Param("id"))
	})

	tests := []struct {
		method string
		path   string
		want   string
	}{
		{"GET", "/users/123", "param: 123"},
		{"POST", "/users/123", "POST user: 123"},
		{"PUT", "/users/123", "PUT user: 123"},
		{"DELETE", "/users/123", "DELETE user: 123"},
		{"PATCH", "/users/123", ""}, // No match
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			handler, params := tree.Find(tt.method, tt.path)
			if tt.want == "" {
				if handler != nil {
					t.Errorf("Find() = %v; want nil", handler)
				}
				return
			}
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			ctx.SetParams(make(map[string]string))
			for _, p := range params {
				ctx.SetParams(map[string]string{p.Key: p.Value})
			}
			handler(w, req, ctx)

			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test complex route scenarios
func TestComplexRouteScenarios(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Add complex routes
	tree.Add("GET", "/api/v1/users/:id/posts/:postId/comments/:commentId", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		userID := ctx.Param("id")
		postID := ctx.Param("postId")
		commentID := ctx.Param("commentId")
		ctx.String(200, "user: %s, post: %s, comment: %s", userID, postID, commentID)
	})

	tree.Add("GET", "/api/v1/users/:id/posts", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		userID := ctx.Param("id")
		ctx.String(200, "user: %s, posts", userID)
	})

	tree.Add("GET", "/api/v1/users", testHandler("users list"))

	tests := []struct {
		path     string
		want     string
		wantCode int
	}{
		{"/api/v1/users", "users list", 200},
		{"/api/v1/users/123/posts", "user: 123, posts", 200},
		{"/api/v1/users/123/posts/456/comments/789", "user: 123, post: 456, comment: 789", 200},
		{"/api/v1/users/123/posts/456", "", 0}, // No match
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, params := tree.Find("GET", tt.path)
			if tt.wantCode == 0 {
				if handler != nil {
					t.Errorf("Find() = %v; want nil", handler)
				}
				return
			}
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			ctx.SetParams(make(map[string]string))
			for _, p := range params {
				ctx.SetParams(map[string]string{p.Key: p.Value})
			}
			handler(w, req, ctx)

			if w.Code != tt.wantCode {
				t.Errorf("Status = %d; want %d", w.Code, tt.wantCode)
			}
			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test edge cases for parameter routes
func TestParameterEdgeCases(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Add parameter routes with edge cases
	tree.Add("GET", "/:id", paramHandler) // Root parameter
	tree.Add("GET", "/users/:id", paramHandler)
	tree.Add("GET", "/users/:id/", paramHandler) // Trailing slash
	tree.Add("GET", "/users/:id/posts/:postId", multiParamHandler)

	tests := []struct {
		path     string
		want     string
		wantCode int
	}{
		{"/123", "param: 123", 200}, // Root parameter
		{"/users/456", "param: 456", 200},
		{"/users/789/", "param: 789", 200}, // Trailing slash
		{"/users/111/posts/222", "user: 111, post: 222", 200},
		{"/", "", 0},      // No match for root
		{"/users", "", 0}, // No match
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, params := tree.Find("GET", tt.path)
			if tt.wantCode == 0 {
				if handler != nil {
					t.Errorf("Find() = %v; want nil", handler)
				}
				return
			}
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			ctx.SetParams(make(map[string]string))
			for _, p := range params {
				ctx.SetParams(map[string]string{p.Key: p.Value})
			}
			handler(w, req, ctx)

			if w.Code != tt.wantCode {
				t.Errorf("Status = %d; want %d", w.Code, tt.wantCode)
			}
			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}

// Test edge cases for wildcard routes
func TestWildcardEdgeCases(t *testing.T) {
	tree := NewRadixTree[Context]()

	// Add wildcard routes with edge cases
	tree.Add("GET", "/*path", wildcardHandler) // Root wildcard
	tree.Add("GET", "/static/*path", wildcardHandler)
	tree.Add("GET", "/files/*path", wildcardHandler)

	tests := []struct {
		path     string
		want     string
		wantCode int
	}{
		{"/anything", "wildcard: anything", 200}, // Root wildcard
		{"/static/css/style.css", "wildcard: css/style.css", 200},
		{"/files/documents/report.pdf", "wildcard: documents/report.pdf", 200},
		{"/", "", 0},       // No match for root
		{"/static", "", 0}, // No match
		{"/files", "", 0},  // No match
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			handler, params := tree.Find("GET", tt.path)
			if tt.wantCode == 0 {
				if handler != nil {
					t.Errorf("Find() = %v; want nil", handler)
				}
				return
			}
			if handler == nil {
				t.Errorf("Find() = nil; want handler")
				return
			}

			// Test the handler
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			ctx := NewContext(req, w)
			ctx.SetParams(make(map[string]string))
			for _, p := range params {
				ctx.SetParams(map[string]string{p.Key: p.Value})
			}
			handler(w, req, ctx)

			if w.Code != tt.wantCode {
				t.Errorf("Status = %d; want %d", w.Code, tt.wantCode)
			}
			if !strings.Contains(w.Body.String(), tt.want) {
				t.Errorf("Response = %s; want to contain %s", w.Body.String(), tt.want)
			}
		})
	}
}
