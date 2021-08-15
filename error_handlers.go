package router

// NotFoundHandler sends a 404 response with the {"error": "Not Found"} message
func NotFoundHandler(c *RouterContext) {
	c.SetStatus(404).JSON(map[string]string{"error": "Page Not found"})
}

// MethodNotAllowed sends a 405 response with the {"error": "Method Not Allowed"} message
func MethodNotAllowed(c *RouterContext) {
	c.SetStatus(405).JSON(map[string]string{"error": "Method Not Allowed"})
}
