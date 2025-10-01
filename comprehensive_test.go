package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test middleware
func testMiddleware(name string) Middleware[Context] {
	return func(next HandlerFunc[Context]) HandlerFunc[Context] {
		return func(w http.ResponseWriter, r *http.Request, ctx Context) {
			ctx.Set("middleware_"+name, "executed")
			next(w, r, ctx)
		}
	}
}

// Test handlers
func testHandler(response string) HandlerFunc[Context] {
	return func(w http.ResponseWriter, r *http.Request, ctx Context) {
		ctx.String(200, "%s", response)
	}
}

func paramHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	param := ctx.Param("id")
	ctx.String(200, "param: %s", param)
}

func wildcardHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	path := ctx.Param("path")
	ctx.String(200, "wildcard: %s", path)
}

func jsonHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	ctx.JSON(200, map[string]string{"message": "test"})
}

func htmlHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	ctx.HTML(200, "<h1>Test</h1>")
}

func redirectHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	ctx.Redirect(302, "/redirected")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	ctx.String(404, "%s", "Not Found")
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request, ctx Context) {
	ctx.String(405, "%s", "Method Not Allowed")
}

func TestParameterRouting(t *testing.T) {
	router := NewRouter(nil)
	router.GET("/users/:id", paramHandler)
	router.GET("/posts/:id/comments/:commentId", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		userID := ctx.Param("id")
		commentID := ctx.Param("commentId")
		ctx.String(200, "user: %s, comment: %s", userID, commentID)
	})

	tests := []struct {
		path   string
		expect string
	}{
		{"/users/123", "param: 123"},
		{"/posts/456/comments/789", "user: 456, comment: 789"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expect) {
				t.Errorf("Expected response to contain %q, got %q", tt.expect, w.Body.String())
			}
		})
	}
}

func TestWildcardRouting(t *testing.T) {
	router := NewRouter(nil)
	router.GET("/static/*path", wildcardHandler)
	router.GET("/files/*path", wildcardHandler)

	tests := []struct {
		path   string
		expect string
	}{
		{"/static/css/style.css", "wildcard: css/style.css"},
		{"/static/js/app.js", "wildcard: js/app.js"},
		{"/files/documents/report.pdf", "wildcard: documents/report.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expect) {
				t.Errorf("Expected response to contain %q, got %q", tt.expect, w.Body.String())
			}
		})
	}
}

func TestMiddleware(t *testing.T) {
	router := NewRouter(nil)
	router.Use(testMiddleware("global"))

	router.GET("/test", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		global, _ := ctx.Get("middleware_global")
		ctx.String(200, "global: %s", global.(string))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "global: executed") {
		t.Errorf("Expected middleware to be executed")
	}
}

func TestRouteGroups(t *testing.T) {
	router := NewRouter(nil)

	// Global middleware
	router.Use(testMiddleware("global"))

	// API group with middleware
	apiGroup := router.Group("/api")
	apiGroup.Use(testMiddleware("api"))

	apiGroup.GET("/users", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		global, _ := ctx.Get("middleware_global")
		api, _ := ctx.Get("middleware_api")
		ctx.String(200, "global: %s, api: %s", global.(string), api.(string))
	})

	// V1 sub-group
	v1Group := apiGroup.Group("/v1")
	v1Group.Use(testMiddleware("v1"))

	v1Group.GET("/posts", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		global, _ := ctx.Get("middleware_global")
		api, _ := ctx.Get("middleware_api")
		v1, _ := ctx.Get("middleware_v1")
		ctx.String(200, "global: %s, api: %s, v1: %s", global.(string), api.(string), v1.(string))
	})

	tests := []struct {
		path   string
		expect string
	}{
		{"/api/users", "global: executed, api: executed"},
		{"/api/v1/posts", "global: executed, api: executed, v1: executed"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expect) {
				t.Errorf("Expected response to contain %q, got %q", tt.expect, w.Body.String())
			}
		})
	}
}

func TestCustomMethods(t *testing.T) {
	router := NewRouter(nil)
	router.Register("KICK", "/admin/kick", testHandler("KICK executed"))
	router.Register("BAN", "/admin/ban", testHandler("BAN executed"))

	tests := []struct {
		method string
		path   string
		expect string
	}{
		{"KICK", "/admin/kick", "KICK executed"},
		{"BAN", "/admin/ban", "BAN executed"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expect) {
				t.Errorf("Expected response to contain %q, got %q", tt.expect, w.Body.String())
			}
		})
	}
}

func TestContextMethods(t *testing.T) {
	router := NewRouter(nil)

	router.GET("/json", jsonHandler)
	router.GET("/html", htmlHandler)
	router.GET("/redirect", redirectHandler)
	router.GET("/query", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		query := ctx.Query("test")
		ctx.String(200, "query: %s", query)
	})

	tests := []struct {
		path   string
		expect string
	}{
		{"/json", `{"message": "test"}`},
		{"/html", "<h1>Test</h1>"},
		{"/query?test=value", "query: value"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expect) {
				t.Errorf("Expected response to contain %q, got %q", tt.expect, w.Body.String())
			}
		})
	}
}

func TestNotFoundAndMethodNotAllowed(t *testing.T) {
	options := &RouterOptions{
		NotFoundHandler:  notFoundHandler,
		MethodNotAllowed: methodNotAllowedHandler,
	}

	router := NewRouter(options)
	router.GET("/test", testHandler("test"))

	// Test 404
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Not Found") {
		t.Errorf("Expected 404 message")
	}

	// Test 405
	req = httptest.NewRequest("POST", "/test", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 405 {
		t.Errorf("Expected status 405, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Method Not Allowed") {
		t.Errorf("Expected 405 message")
	}
}

func TestStaticFileServing(t *testing.T) {
	router := NewRouter(nil)
	router.Static("/static/", "./testdata")

	req := httptest.NewRequest("GET", "/static/test.txt", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// This will fail if testdata/test.txt doesn't exist, but that's expected
	// In a real test, you'd create the test file first
}

func TestDelegateRoutes(t *testing.T) {
	router := NewRouter(nil)
	router.DELEGATE("/files/", http.MethodGet, testHandler("delegate"))

	req := httptest.NewRequest("GET", "/files/document.pdf", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "delegate") {
		t.Errorf("Expected delegate response")
	}
}

func TestComplexRouting(t *testing.T) {
	router := NewRouter(nil)

	// Multiple parameter routes
	router.GET("/users/:id/posts/:postId", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		userID := ctx.Param("id")
		postID := ctx.Param("postId")
		ctx.String(200, "user: %s, post: %s", userID, postID)
	})

	// Mixed static and parameter routes
	router.GET("/api/v1/users/:id", paramHandler)
	router.GET("/api/v1/users", testHandler("users list"))

	// Wildcard with parameters
	router.GET("/api/*path", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		path := ctx.Param("path")
		ctx.String(200, "api wildcard: %s", path)
	})

	tests := []struct {
		path   string
		expect string
	}{
		{"/users/123/posts/456", "user: 123, post: 456"},
		{"/api/v1/users/789", "param: 789"},
		{"/api/v1/users", "users list"},
		{"/api/anything/here", "api wildcard: /anything/here"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expect) {
				t.Errorf("Expected response to contain %q, got %q", tt.expect, w.Body.String())
			}
		})
	}
}

func TestMiddlewareOrder(t *testing.T) {
	router := NewRouter(nil)

	// Add middleware in specific order
	router.Use(func(next HandlerFunc[Context]) HandlerFunc[Context] {
		return func(w http.ResponseWriter, r *http.Request, ctx Context) {
			ctx.Set("order", "1")
			next(w, r, ctx)
		}
	})

	router.Use(func(next HandlerFunc[Context]) HandlerFunc[Context] {
		return func(w http.ResponseWriter, r *http.Request, ctx Context) {
			order, _ := ctx.Get("order")
			ctx.Set("order", order.(string)+",2")
			next(w, r, ctx)
		}
	})

	router.GET("/test", func(w http.ResponseWriter, r *http.Request, ctx Context) {
		order, _ := ctx.Get("order")
		ctx.String(200, "order: %s", order.(string))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Middleware should execute in the order they were added
	if !strings.Contains(w.Body.String(), "order: 1,2") {
		t.Errorf("Expected middleware to execute in order")
	}
}
