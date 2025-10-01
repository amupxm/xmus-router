package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	router "github.com/amupxm/xmus-router"
)

// Middleware functions
func loggingMiddleware(next router.HandlerFunc[router.Context]) router.HandlerFunc[router.Context] {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		start := time.Now()
		fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)

		next(w, r, ctx)

		duration := time.Since(start)
		fmt.Printf("Response time: %v\n", duration)
	}
}

func authMiddleware(next router.HandlerFunc[router.Context]) router.HandlerFunc[router.Context] {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		// Simple auth check
		auth := r.Header.Get("Authorization")
		if auth == "" {
			ctx.String(401, "Unauthorized")
			return
		}

		// Set user info in context
		ctx.Set("user", "authenticated_user")
		next(w, r, ctx)
	}
}

func corsMiddleware(next router.HandlerFunc[router.Context]) router.HandlerFunc[router.Context] {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			ctx.String(200, "OK")
			return
		}

		next(w, r, ctx)
	}
}

// Handler functions
func homeHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	ctx.HTML(200, "<h1>Welcome to XMUS Router!</h1><p>Fast, lightweight router with middleware support</p>")
}

func userHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	userID := ctx.Param("id")
	ctx.JSON(200, map[string]string{
		"message": "User profile",
		"user_id": userID,
	})
}

func createUserHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	ctx.JSON(201, map[string]string{
		"message": "User created successfully",
	})
}

func updateUserHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	userID := ctx.Param("id")
	ctx.JSON(200, map[string]string{
		"message": "User updated",
		"user_id": userID,
	})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	userID := ctx.Param("id")
	ctx.JSON(200, map[string]string{
		"message": "User deleted",
		"user_id": userID,
	})
}

func postsHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	ctx.JSON(200, map[string]interface{}{
		"posts": []map[string]string{
			{"id": "1", "title": "First Post"},
			{"id": "2", "title": "Second Post"},
		},
	})
}

func postHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	postID := ctx.Param("id")
	ctx.JSON(200, map[string]string{
		"post_id": postID,
		"title":   "Sample Post",
		"content": "This is a sample post content",
	})
}

func wildcardHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	path := ctx.Param("path")
	ctx.String(200, "Wildcard route matched: %s", path)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	ctx.String(404, "Custom 404: Page not found")
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request, ctx router.Context) {
	ctx.String(405, "Method not allowed")
}

func main() {
	// Create router with options
	options := &router.RouterOptions{
		NotFoundHandler:  notFoundHandler,
		MethodNotAllowed: methodNotAllowedHandler,
	}

	rt := router.NewRouter(options)

	// Global middleware
	rt.Use(loggingMiddleware, corsMiddleware)

	// Basic routes
	rt.GET("/", homeHandler)
	rt.GET("/wildcard/*path", wildcardHandler)

	// User routes with parameters
	rt.GET("/users/:id", userHandler)
	rt.POST("/users", createUserHandler)
	rt.PUT("/users/:id", updateUserHandler)
	rt.DELETE("/users/:id", deleteUserHandler)

	// Route groups with middleware
	apiGroup := rt.Group("/api")
	apiGroup.Use(authMiddleware)

	// API v1 group
	v1Group := apiGroup.Group("/v1")
	v1Group.GET("/posts", postsHandler)
	v1Group.GET("/posts/:id", postHandler)
	v1Group.POST("/posts", func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		ctx.JSON(201, map[string]string{"message": "Post created"})
	})

	// API v2 group with different middleware
	v2Group := apiGroup.Group("/v2")
	v2Group.GET("/posts", postsHandler)

	// Static file serving
	rt.Static("/static/", "./static")

	// Custom method registration
	rt.Register("KICK", "/admin/kick", func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		ctx.String(200, "Custom KICK method executed")
	})

	// Delegate route for file serving
	rt.DELEGATE("/files/", http.MethodGet, func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		ctx.String(200, "File serving delegate")
	})

	fmt.Println("Server starting on :8080")
	fmt.Println("Available routes:")
	fmt.Println("  GET  /")
	fmt.Println("  GET  /users/:id")
	fmt.Println("  POST /users")
	fmt.Println("  PUT  /users/:id")
	fmt.Println("  DELETE /users/:id")
	fmt.Println("  GET  /api/v1/posts")
	fmt.Println("  GET  /api/v1/posts/:id")
	fmt.Println("  POST /api/v1/posts")
	fmt.Println("  GET  /api/v2/posts")
	fmt.Println("  GET  /wildcard/*path")
	fmt.Println("  KICK /admin/kick")
	fmt.Println("  GET  /static/*")
	fmt.Println("  GET  /files/*")

	log.Fatal(http.ListenAndServe(":8080", rt))
}
