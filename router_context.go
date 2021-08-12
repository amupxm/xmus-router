package router

import "encoding/json"

func (rc *RouterContext) JSON(data interface{}) {
	rc.SetHeader("Content-Type", "application/json")
	d, _ := json.Marshal(data)
	rc.Response.Write(d)
}
func (rc *RouterContext) SetHeader(key string, value string) *RouterContext {
	rc.Response.Header().Set(key, value)
	return rc
}
func (rc *RouterContext) SetStatus(status int) *RouterContext {
	rc.Response.WriteHeader(status)
	return rc
}
