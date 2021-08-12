package logger

import (
	"fmt"
	"log"

	router "github.com/amupxm/xmus-router"
)

var (
	greenBg   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	whiteBg   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellowBg  = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	redBg     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blueBg    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magentaBg = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyanBg    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	green     = string([]byte{27, 91, 51, 50, 109})
	white     = string([]byte{27, 91, 51, 55, 109})
	yellow    = string([]byte{27, 91, 51, 51, 109})
	red       = string([]byte{27, 91, 51, 49, 109})
	blue      = string([]byte{27, 91, 51, 52, 109})
	magenta   = string([]byte{27, 91, 51, 53, 109})
	cyan      = string([]byte{27, 91, 51, 54, 109})
	reset     = string([]byte{27, 91, 48, 109})
)

func Logger(c *router.RouterContext) {
	b := whiteBg
	switch c.Request.Method {
	case "GET":
		b = greenBg
	case "POST":
		b = blueBg
	case "PUT":
		b = yellowBg
	case "DELETE":
		b = redBg
	case "PATCH":
		b = yellowBg
	default:
		b = whiteBg
	}

	log.Printf("%s | %s | %s ", colorPrint(c.Request.Method, b), c.Request.URL.Path, c.Request.Proto)
}
func colorPrint(format string, color string, s ...interface{}) string {
	if len(s) == 0 {
		return fmt.Sprint(color + format + reset)
	}
	data := fmt.Sprintf(format, s...)
	return fmt.Sprint(color + data + reset)
}
