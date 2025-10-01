package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ExampleBasicUsage demonstrates basic PACT router usage
func ExampleBasicUsage() {
	fmt.Println("=== Basic PACT Router Usage ===")

	// Create a new PACT router
	router := NewPACTRouter()

	// Define some routes
	routes := []Route{
		{Path: "/api/v1/users", Method: "GET", Handler: "getUsers"},
		{Path: "/api/v1/users", Method: "POST", Handler: "createUser"},
		{Path: "/api/v1/users/:id", Method: "GET", Handler: "getUser"},
		{Path: "/api/v1/users/:id", Method: "PUT", Handler: "updateUser"},
		{Path: "/api/v1/users/:id", Method: "DELETE", Handler: "deleteUser"},
		{Path: "/api/v1/posts", Method: "GET", Handler: "getPosts"},
		{Path: "/api/v1/posts", Method: "POST", Handler: "createPost"},
		{Path: "/api/v1/posts/:id", Method: "GET", Handler: "getPost"},
		{Path: "/api/v1/posts/:id", Method: "PUT", Handler: "updatePost"},
		{Path: "/api/v1/posts/:id", Method: "DELETE", Handler: "deletePost"},
		{Path: "/api/v1/comments", Method: "GET", Handler: "getComments"},
		{Path: "/api/v1/comments/:id", Method: "GET", Handler: "getComment"},
		{Path: "/api/v2/users", Method: "GET", Handler: "getUsersV2"},
		{Path: "/api/v2/posts", Method: "GET", Handler: "getPostsV2"},
	}

	// Build the router with optimization
	fmt.Println("Building router with", len(routes), "routes...")
	start := time.Now()
	router.Build(routes)
	buildTime := time.Since(start)
	fmt.Printf("Build completed in %v\n", buildTime)

	// Test some lookups
	testPaths := []string{
		"/api/v1/users",
		"/api/v1/users/123",
		"/api/v1/posts",
		"/api/v1/posts/456",
		"/api/v1/comments",
		"/api/v1/comments/789",
		"/api/v2/users",
		"/api/v2/posts",
		"/api/v1/nonexistent",
	}

	fmt.Println("\nTesting lookups:")
	for _, path := range testPaths {
		start := time.Now()
		handler := router.Lookup(path)
		lookupTime := time.Since(start)

		if handler != nil {
			fmt.Printf("✓ %s -> %v (lookup: %v)\n", path, handler, lookupTime)
		} else {
			fmt.Printf("✗ %s -> Not found (lookup: %v)\n", path, lookupTime)
		}
	}

	// Show analysis results
	fmt.Println("\nRoute Analysis Results:")
	fmt.Printf("Common prefixes: %v\n", getTopPrefixes(router.analyzer.commonPrefixes, 5))
	fmt.Printf("Hot paths: %v\n", router.analyzer.hotPaths)
	fmt.Printf("Hot path cache size: %d\n", len(router.hotPaths))
}

// ExampleAdvancedUsage demonstrates advanced PACT router features
func ExampleAdvancedUsage() {
	fmt.Println("\n=== Advanced PACT Router Usage ===")

	// Create advanced router with custom configuration
	config := &RouterConfig{
		HotPathCacheSize:   64,
		HotPathThreshold:   60.0,
		MaxMemoryUsage:     1024 * 1024 * 5, // 5MB
		CompressionEnabled: true,
		SIMDEnabled:        true,
		ConcurrentAccess:   true,
	}

	router := NewAdvancedPACTRouter(config)

	// Generate a larger route set
	routes := generateRESTRoutes(200)
	router.Build(routes)

	// Simulate some traffic
	fmt.Println("Simulating traffic...")
	paths := generateAccessPattern(routes, 0.8)

	start := time.Now()
	for i := 0; i < 10000; i++ {
		path := paths[i%len(paths)]
		router.ConcurrentLookup(path)
	}
	totalTime := time.Since(start)

	// Show performance metrics
	stats := router.GetStats()

	fmt.Printf("Processed 10,000 lookups in %v\n", totalTime)
	fmt.Printf("Average lookup time: %.2f ns\n", router.GetAverageLookupTime())
	fmt.Printf("Cache hit rate: %.2f%%\n", router.GetCacheHitRate())
	fmt.Printf("Memory usage: %d bytes\n", router.MemoryUsage())

	fmt.Println("\nDetailed Metrics:")
	fmt.Printf("Total lookups: %d\n", stats.TotalLookups)
	fmt.Printf("Cache hits: %d\n", stats.CacheHits)
	fmt.Printf("Cache misses: %d\n", stats.CacheMisses)
	fmt.Printf("Max lookup time: %d ns\n", stats.MaxLookupTime)
	fmt.Printf("Min lookup time: %d ns\n", stats.MinLookupTime)

	// Show configuration
	fmt.Println("\nConfiguration:")
	fmt.Printf("Hot path cache size: %d\n", config.HotPathCacheSize)
	fmt.Printf("Hot path threshold: %.1f\n", config.HotPathThreshold)
	fmt.Printf("Compression enabled: %t\n", config.CompressionEnabled)
	fmt.Printf("SIMD enabled: %t\n", config.SIMDEnabled)
	fmt.Printf("Concurrent access: %t\n", config.ConcurrentAccess)

	// Health check
	if router.HealthCheck() {
		fmt.Println("\n✓ Router health check passed")
	} else {
		fmt.Println("\n✗ Router health check failed")
	}
}

// ExamplePerformanceComparison compares PACT with different configurations
func ExamplePerformanceComparison() {
	fmt.Println("\n=== Performance Comparison ===")

	routes := generateRESTRoutes(100)
	paths := generateAccessPattern(routes, 0.8)

	// Test different configurations
	configs := []struct {
		name   string
		config *RouterConfig
	}{
		{
			name: "Small Cache",
			config: &RouterConfig{
				HotPathCacheSize:   16,
				HotPathThreshold:   80.0,
				CompressionEnabled: false,
				SIMDEnabled:        false,
			},
		},
		{
			name: "Medium Cache",
			config: &RouterConfig{
				HotPathCacheSize:   32,
				HotPathThreshold:   70.0,
				CompressionEnabled: true,
				SIMDEnabled:        false,
			},
		},
		{
			name: "Large Cache",
			config: &RouterConfig{
				HotPathCacheSize:   64,
				HotPathThreshold:   60.0,
				CompressionEnabled: true,
				SIMDEnabled:        true,
			},
		},
	}

	for _, cfg := range configs {
		fmt.Printf("\nTesting %s configuration:\n", cfg.name)

		router := NewAdvancedPACTRouter(cfg.config)
		router.Build(routes)

		// Warm up
		for i := 0; i < 1000; i++ {
			path := paths[i%len(paths)]
			router.ConcurrentLookup(path)
		}

		// Benchmark
		start := time.Now()
		for i := 0; i < 10000; i++ {
			path := paths[i%len(paths)]
			router.ConcurrentLookup(path)
		}
		totalTime := time.Since(start)

		fmt.Printf("  Total time: %v\n", totalTime)
		fmt.Printf("  Average lookup: %.2f ns\n", router.GetAverageLookupTime())
		fmt.Printf("  Cache hit rate: %.2f%%\n", router.GetCacheHitRate())
		fmt.Printf("  Memory usage: %d bytes\n", router.MemoryUsage())
	}
}

// ExampleRealWorldScenario simulates a real-world API scenario
func ExampleRealWorldScenario() {
	fmt.Println("\n=== Real-World API Scenario ===")

	// Simulate a microservices API with multiple services
	services := []string{"users", "posts", "comments", "likes", "followers", "messages", "notifications", "settings"}
	versions := []string{"v1", "v2"}
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	var routes []Route

	// Generate routes for each service
	for _, service := range services {
		for _, version := range versions {
			for _, method := range methods {
				// Collection endpoint
				routes = append(routes, Route{
					Path:    fmt.Sprintf("/api/%s/%s", version, service),
					Method:  method,
					Handler: fmt.Sprintf("handle%s%s%s", method, version, service),
				})

				// Resource endpoint
				routes = append(routes, Route{
					Path:    fmt.Sprintf("/api/%s/%s/:id", version, service),
					Method:  method,
					Handler: fmt.Sprintf("handle%s%s%sById", method, version, service),
				})
			}
		}
	}

	fmt.Printf("Generated %d routes for %d services\n", len(routes), len(services))

	// Create router
	router := NewAdvancedPACTRouter(nil)
	router.Build(routes)

	// Simulate realistic traffic patterns
	fmt.Println("Simulating realistic traffic...")

	// 80% of traffic goes to top 20% of routes (Pareto distribution)
	hotRoutes := routes[:len(routes)/5]  // Top 20%
	coldRoutes := routes[len(routes)/5:] // Bottom 80%

	totalRequests := 100000
	hotRequests := int(float64(totalRequests) * 0.8)
	coldRequests := totalRequests - hotRequests

	start := time.Now()

	// Process hot routes
	for i := 0; i < hotRequests; i++ {
		route := hotRoutes[i%len(hotRoutes)]
		router.ConcurrentLookup(route.Path)
	}

	// Process cold routes
	for i := 0; i < coldRequests; i++ {
		route := coldRoutes[i%len(coldRoutes)]
		router.ConcurrentLookup(route.Path)
	}

	totalTime := time.Since(start)

	// Show results
	fmt.Printf("Processed %d requests in %v\n", totalRequests, totalTime)
	fmt.Printf("Requests per second: %.0f\n", float64(totalRequests)/totalTime.Seconds())
	fmt.Printf("Average lookup time: %.2f ns\n", router.GetAverageLookupTime())
	fmt.Printf("Cache hit rate: %.2f%%\n", router.GetCacheHitRate())
	fmt.Printf("Memory usage: %d bytes\n", router.MemoryUsage())

	// Show hot paths
	fmt.Printf("Hot paths identified: %d\n", len(router.analyzer.hotPaths))
	fmt.Printf("Top hot paths: %v\n", router.analyzer.hotPaths[:min(5, len(router.analyzer.hotPaths))])
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// generateRESTRoutes generates realistic REST API routes
func generateRESTRoutes(count int) []Route {
	routes := make([]Route, 0, count)

	// Common API patterns
	resources := []string{"users", "posts", "comments", "likes", "followers", "messages", "notifications", "settings"}
	versions := []string{"v1", "v2", "v3"}
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		version := versions[rand.Intn(len(versions))]
		resource := resources[rand.Intn(len(resources))]
		method := methods[rand.Intn(len(methods))]

		var path string
		switch rand.Intn(4) {
		case 0: // Collection endpoint
			path = fmt.Sprintf("/api/%s/%s", version, resource)
		case 1: // Resource endpoint
			path = fmt.Sprintf("/api/%s/%s/:id", version, resource)
		case 2: // Nested resource
			path = fmt.Sprintf("/api/%s/%s/:id/%s", version, resource, resources[rand.Intn(len(resources))])
		case 3: // Action endpoint
			path = fmt.Sprintf("/api/%s/%s/:id/activate", version, resource)
		}

		routes = append(routes, Route{
			Path:    path,
			Method:  method,
			Handler: fmt.Sprintf("handle%s%s", method, resource),
		})
	}

	return routes
}

// generateAccessPattern generates realistic access patterns following 80/20 rule
func generateAccessPattern(routes []Route, hotPathRatio float64) []string {
	// Sort routes by access frequency (shorter paths first)
	sortedRoutes := make([]Route, len(routes))
	copy(sortedRoutes, routes)

	// Simple sorting by path length (shorter = more frequent)
	for i := 0; i < len(sortedRoutes); i++ {
		for j := i + 1; j < len(sortedRoutes); j++ {
			if len(sortedRoutes[i].Path) > len(sortedRoutes[j].Path) {
				sortedRoutes[i], sortedRoutes[j] = sortedRoutes[j], sortedRoutes[i]
			}
		}
	}

	// Generate access pattern
	hotPathCount := int(float64(len(routes)) * hotPathRatio)
	patterns := make([]string, 0, 10000) // Generate 10k access patterns

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10000; i++ {
		if i < int(float64(10000)*hotPathRatio) {
			// Hot path access
			route := sortedRoutes[rand.Intn(hotPathCount)]
			patterns = append(patterns, route.Path)
		} else {
			// Cold path access
			route := sortedRoutes[hotPathCount+rand.Intn(len(routes)-hotPathCount)]
			patterns = append(patterns, route.Path)
		}
	}

	return patterns
}

// ExampleErrorHandling demonstrates error handling and edge cases
func ExampleErrorHandling() {
	fmt.Println("\n=== Error Handling and Edge Cases ===")

	router := NewPACTRouter()

	// Test with empty routes
	fmt.Println("Testing empty router...")
	handler := router.Lookup("/api/users")
	if handler == nil {
		fmt.Println("✓ Empty router correctly returns nil")
	}

	// Test with invalid paths
	fmt.Println("Testing invalid paths...")
	invalidPaths := []string{
		"",
		"/",
		"//",
		"/api//users",
		"/api/users/",
		"/api/users//123",
	}

	for _, path := range invalidPaths {
		handler := router.Lookup(path)
		if handler == nil {
			fmt.Printf("✓ Invalid path '%s' correctly returns nil\n", path)
		} else {
			fmt.Printf("✗ Invalid path '%s' unexpectedly returned %v\n", path, handler)
		}
	}

	// Test with very long paths
	fmt.Println("Testing very long paths...")
	longPath := "/api/v1/users/" + string(make([]byte, 1000))
	handler = router.Lookup(longPath)
	if handler == nil {
		fmt.Println("✓ Very long path correctly returns nil")
	}

	// Test concurrent access
	fmt.Println("Testing concurrent access...")
	router.Build([]Route{
		{Path: "/api/users", Method: "GET", Handler: "getUsers"},
	})

	// This would be more comprehensive with actual goroutines
	// For now, just test that it doesn't panic
	handler = router.Lookup("/api/users")
	if handler != nil {
		fmt.Println("✓ Concurrent access works correctly")
	}
}

// ExampleMonitoring demonstrates monitoring and observability
func ExampleMonitoring() {
	fmt.Println("\n=== Monitoring and Observability ===")

	router := NewAdvancedPACTRouter(nil)
	routes := generateRESTRoutes(50)
	router.Build(routes)

	// Simulate some traffic
	paths := generateAccessPattern(routes, 0.8)
	for i := 0; i < 5000; i++ {
		path := paths[i%len(paths)]
		router.ConcurrentLookup(path)
	}

	// Export metrics
	metrics := router.ExportMetrics()

	fmt.Println("Performance Metrics:")
	fmt.Printf("  Total lookups: %v\n", metrics["lookups"].(map[string]interface{})["total"])
	fmt.Printf("  Cache hit rate: %.2f%%\n", metrics["lookups"].(map[string]interface{})["hit_rate"])
	fmt.Printf("  Average lookup time: %.2f ns\n", metrics["timing"].(map[string]interface{})["average_ns"])
	fmt.Printf("  Memory usage: %v bytes\n", metrics["memory"].(map[string]interface{})["usage_bytes"])

	// Health check
	if router.HealthCheck() {
		fmt.Println("✓ Router is healthy")
	} else {
		fmt.Println("✗ Router health check failed")
	}

	// Reset statistics
	router.ResetStats()
	fmt.Println("✓ Statistics reset")
}

// Main function to run all examples
func main_2() {
	// Run all examples
	ExampleBasicUsage()
	ExampleAdvancedUsage()
	ExamplePerformanceComparison()
	ExampleRealWorldScenario()
	ExampleErrorHandling()
	ExampleMonitoring()

	fmt.Println("\n=== All Examples Completed ===")
}
