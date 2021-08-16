package router

import (
	"errors"
	"regexp"
)

var error404 = []byte(`{"error": "Page Not found"}`)
var error405 = []byte(`{"error": "Method Not Allowed"}`)

var validPathStartAndEndRegex = regexp.MustCompile(`^\/(.?)*\/$`)

var validateRequestPathRegex = regexp.MustCompile(`^(.?)*\/$`)

var delegateRegex, _ = regexp.Compile(`\/((.?)*\/)*\*\/$`)

var hasParamsRegex, _ = regexp.Compile(`\:(\w|\d)*`)

var getURLParamsRegex = regexp.MustCompile(`\:((\W^\/)|\w|\d)*\/`)

var validatePathRegex, _ = regexp.Compile(`\/((.?)*\/)*`)

var hasSpaceRegex, _ = regexp.Compile(`\s`)

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

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
	MethodPatch  = "PATCH"
)

var errMethodNotAllowed = errors.New("405")
var errNotFound = errors.New("404")

var errorNotFoundMessage = []byte(`{"error":"Not found"}`)
var errorMethodNotAllowedMessage = []byte(`{"error":"method not allowed"}`)
