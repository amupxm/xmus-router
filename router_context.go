package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (rc XmusContext) JSON(data interface{}) {
	rc.SetHeader("Content-Type", "application/json")
	d, _ := json.Marshal(data)
	rc.Response.Write(d)
}
func (rc XmusContext) SetHeader(key string, value string) XmusContext {
	rc.Response.Header().Set(key, value)
	return rc
}
func (rc XmusContext) SetStatus(status int) XmusContext {
	rc.Response.WriteHeader(status)
	return rc
}

type (
	XmusContext struct {
		Response  http.ResponseWriter
		Request   *http.Request
		URLParams map[string]string
	}
)

func (xm *XmusContext) buildParams(requestPath, path []string) {
	log.Println(requestPath, path)
	params := make(map[string]string)
	for i, v := range path {
		if strings.HasPrefix(v, ":") {
			params[path[i][1:]] = requestPath[i]
		}
	}
	xm.URLParams = params
}
func (xm XmusContext) GetParam(key string) string {
	return xm.URLParams[key]
}
