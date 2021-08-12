package router

func NotFoundHandler(c *RouterContext) {
	c.SetStatus(404).JSON(map[string]string{"error": "Page Not found"})
}

func MethodNotAllowed(c *RouterContext) {
	c.SetStatus(405).JSON(map[string]string{"error": "Method Not Allowed"})
}
