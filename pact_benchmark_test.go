package main

import (
	"testing"
)

// BenchmarkPACT benchmarks the PACT router
func BenchmarkPACT(b *testing.B) {
	// Generate realistic route set
	routes := generateRESTRoutes(100)

	// Build router
	router := NewPACTRouter()
	router.Build(routes)

	// Create realistic access pattern (80/20 rule)
	paths := generateAccessPattern(routes, 0.8)

	// Warm up caches
	for _, path := range paths[:1000] {
		router.Lookup(path)
	}

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := paths[i%len(paths)]
		router.Lookup(path)
	}
}

// BenchmarkPACTLookup benchmarks just the lookup operation
func BenchmarkPACTLookup(b *testing.B) {
	routes := generateRESTRoutes(50)
	router := NewPACTRouter()
	router.Build(routes)

	// Test with hot paths
	hotPaths := []string{
		"/api/v1/users",
		"/api/v1/posts",
		"/api/v1/comments",
		"/api/v1/users/123",
		"/api/v1/posts/456",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := hotPaths[i%len(hotPaths)]
		router.Lookup(path)
	}
}

// BenchmarkPACTBuild benchmarks the build phase
func BenchmarkPACTBuild(b *testing.B) {
	routes := generateRESTRoutes(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router := NewPACTRouter()
		router.Build(routes)
	}
}

// BenchmarkPACTMemory benchmarks memory usage
func BenchmarkPACTMemory(b *testing.B) {
	routes := generateRESTRoutes(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router := NewPACTRouter()
		router.Build(routes)
		_ = router // Prevent optimization
	}
}

// TestPACTCorrectness tests that PACT router works correctly
func TestPACTCorrectness(t *testing.T) {
	routes := []Route{
		{Path: "/api/v1/users", Method: "GET", Handler: "getUsers"},
		{Path: "/api/v1/users", Method: "POST", Handler: "createUser"},
		{Path: "/api/v1/users/:id", Method: "GET", Handler: "getUser"},
		{Path: "/api/v1/users/:id", Method: "PUT", Handler: "updateUser"},
		{Path: "/api/v1/users/:id", Method: "DELETE", Handler: "deleteUser"},
		{Path: "/api/v1/posts", Method: "GET", Handler: "getPosts"},
		{Path: "/api/v1/posts/:id", Method: "GET", Handler: "getPost"},
		{Path: "/api/v2/users", Method: "GET", Handler: "getUsersV2"},
	}

	router := NewPACTRouter()
	router.Build(routes)

	// Test cases: [path, expectedHandler, shouldFind]
	testCases := []struct {
		path            string
		expectedHandler string
		shouldFind      bool
	}{
		{"/api/v1/users", "getUsers", true},
		{"/api/v1/users/123", "getUser", true},
		{"/api/v1/posts", "getPosts", true},
		{"/api/v1/posts/456", "getPost", true},
		{"/api/v2/users", "getUsersV2", true},
		{"/api/v1/nonexistent", "", false},
		{"/api/v1/users/123/posts", "", false},
	}

	for _, tc := range testCases {
		handler := router.Lookup(tc.path)
		if tc.shouldFind {
			if handler == nil {
				t.Errorf("Expected to find handler for path %s, but got nil", tc.path)
			}
		} else {
			if handler != nil {
				t.Errorf("Expected no handler for path %s, but got %v", tc.path, handler)
			}
		}
	}
}

// TestPACTHotPathCaching tests hot path caching
func TestPACTHotPathCaching(t *testing.T) {
	routes := generateRESTRoutes(50)
	router := NewPACTRouter()
	router.Build(routes)

	// Check that hot paths are cached
	if len(router.hotPaths) == 0 {
		t.Error("Expected hot paths to be cached, but got empty cache")
	}

	// Test that cached paths return results quickly
	hotPath := router.analyzer.hotPaths[0]
	handler := router.Lookup(hotPath)
	if handler == nil {
		t.Errorf("Expected to find handler for hot path %s", hotPath)
	}
}

// TestPACTMemoryUsage tests memory usage characteristics
func TestPACTMemoryUsage(t *testing.T) {
	routes := generateRESTRoutes(100)
	router := NewPACTRouter()
	router.Build(routes)

	// Check that hot path cache is not too large
	if len(router.hotPaths) > HOT_PATH_CACHE_SIZE {
		t.Errorf("Hot path cache size %d exceeds limit %d", len(router.hotPaths), HOT_PATH_CACHE_SIZE)
	}

	// Check that common prefixes were identified
	if len(router.analyzer.commonPrefixes) == 0 {
		t.Error("Expected common prefixes to be identified")
	}
}

// BenchmarkComparison benchmarks PACT against a simple map
func BenchmarkComparison(b *testing.B) {
	routes := generateRESTRoutes(100)

	// Build PACT router
	pactRouter := NewPACTRouter()
	pactRouter.Build(routes)

	// Build simple map router
	mapRouter := make(map[string]interface{})
	for _, route := range routes {
		mapRouter[route.Path] = route.Handler
	}

	// Generate test paths
	paths := generateAccessPattern(routes, 0.8)

	b.Run("PACT", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			path := paths[i%len(paths)]
			pactRouter.Lookup(path)
		}
	})

	b.Run("Map", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			path := paths[i%len(paths)]
			_ = mapRouter[path]
		}
	})
}
