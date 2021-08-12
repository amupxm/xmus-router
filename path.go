package router

import (
	"fmt"
	"strings"
)

func (rt *router) PrepareURLParams(path string) (map[string]string, error) {
	// url can be contain /param/:{param}/:regex
	// split path by /
	params := make(map[string]string)
	paths := strings.Split(path, "/")
	for i := range paths {
		if strings.HasPrefix(paths[i], ":") {

			if _, ok := params[paths[i][1:]]; ok {
				return nil, fmt.Errorf("duplicated %s key in path", paths[i][1:])
			}
			params[paths[i][1:]] = ""
		}
	}
	return params, nil
}
func (rt *router) extractUrlParams(requested, exact string) (map[string]string, error) {
	if !strings.HasSuffix(requested, "/") {
		requested = fmt.Sprintf("%s/", requested)
	}
	result := make(map[string]string)
	rArr := strings.Split(requested, "/")
	eArr := strings.Split(exact, "/")
	for i := range rArr {
		if strings.HasPrefix(eArr[i], ":") {
			result[eArr[i][1:]] = rArr[i]
		}
	}
	return result, nil
}
func (rt *router) isMatchedPath(routePath, requestedPath string) bool {
	if !strings.HasSuffix(requestedPath, "/") {
		requestedPath = fmt.Sprintf("%s/", requestedPath)
	}
	rpArr := strings.Split(requestedPath, "/")
	epArr := strings.Split(routePath, "/")
	if len(rpArr) != len(epArr) {
		return false
	}
	for i := range rpArr {
		if !strings.HasPrefix(epArr[i], ":") && epArr[i] != rpArr[i] {
			return false
		}
	}
	return true
}
