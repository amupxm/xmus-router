package connection_utils

import router "github.com/amupxm/xmus-router"

func UpgraderXMethod(c *router.RouterContext) {
	xMethod := c.Request.Header.Get("X-Method")
	if xMethod != "" {
		c.Request.Method = xMethod
	}
}
func LimitRequestBody(c *router.RouterContext, maxBodySize int64) {
	if c.Request.ContentLength > maxBodySize {
		c.SetStatus(413).JSON(`{"error":"Request body is too large"}`)
		r := c.Request.Context()
		r.Done()
		return
	}
}

func CancelAll(c *router.RouterContext, maxBodySize int64) {
	c.SetStatus(413).JSON(`{"error":"calcelled"}`)
	r := c.Request.Context()
	r.Done()
	return
}
