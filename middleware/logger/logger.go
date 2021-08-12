package logger

import (
	"log"

	router "github.com/amupxm/xmus-router"
)

func Logger(c *router.RouterContext) {
	log.Printf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto)
}
