package main

import (
	"fmt"
	"time"
)

// Demo demonstrates the PACT router functionality
func main() {
	fmt.Println("=== PACT Router Demo ===")

	// Create router
	router := NewPACTRouter()

	// Define routes
	routes := []Route{
		{Path: "/api/v1/users", Method: "GET", Handler: "getUsers"},
		{Path: "/api/v1/users", Method: "POST", Handler: "createUser"},
		{Path: "/api/v1/users/:id", Method: "GET", Handler: "getUser"},
		{Path: "/api/v1/users/:id", Method: "PUT", Handler: "updateUser"},
		{Path: "/api/v1/users/:id", Method: "DELETE", Handler: "deleteUser"},
		{Path: "/api/v1/posts", Method: "GET", Handler: "getPosts"},
		{Path: "/api/v1/posts", Method: "POST", Handler: "createPost"},
		{Path: "/api/v1/posts/:id", Method: "GET", Handler: "getPost"},
		{Path: "/api/v2/users", Method: "GET", Handler: "getUsersV2"},
	}

	// Build router with optimization
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
		"/api/v2/users",
		"/api/v1/nonexistent",
		"/api/v1/users/123/posts",
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
	fmt.Printf("Common prefixes: %v\n", getTopPrefixes(router.analyzer.commonPrefixes, 3))
	fmt.Printf("Hot paths: %v\n", router.analyzer.hotPaths)
	fmt.Printf("Hot path cache size: %d\n", len(router.hotPaths))

	// Performance test
	fmt.Println("\nPerformance Test:")
	iterations := 1000000
	start = time.Now()
	for i := 0; i < iterations; i++ {
		router.Lookup("/api/v1/users")
	}
	totalTime := time.Since(start)
	avgTime := totalTime / time.Duration(iterations)

	fmt.Printf("Performed %d lookups in %v\n", iterations, totalTime)
	fmt.Printf("Average lookup time: %v\n", avgTime)
	fmt.Printf("Lookups per second: %.0f\n", float64(iterations)/totalTime.Seconds())

	fmt.Println("\n=== Demo Completed ===")
}
