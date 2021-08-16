package router

import (
	"fmt"
	"strings"
)

func validatePath(path string) string {

	if strings.Contains(path, "//") {
		panic("p[ath must not inclide //")
	}
	if path == "" {
		path = "/"
	}
	if path != "/" {
		if ok := validPathStartAndEndRegex.MatchString(path); !ok {
			panic("path must start and end with /")
		}
	}
	URLParams := getURLParamsRegex.FindAllString(path, -1)
	pathArr := strings.Split(path, "/")
	for _, p := range URLParams {
		if ok := isParamKey(pathArr, p); ok {
			panic("param key is duplicated")
		}
	}
	return path
}

func prepareRequestPath(path string) string {
	if path == "" {
		path = "/"
	}
	if path != "/" && len(path) > 1 {
		if !validateRequestPathRegex.MatchString(path) {
			path = fmt.Sprintf("%s/", path)
		}
	}
	return path
}

func getPathInfo(path string) (hasParams, isDelegate bool, URLParams []string) {
	isDelegate = delegateRegex.MatchString(path)
	hasParams = hasParamsRegex.MatchString(path)
	if hasParams {
		URLParams = getURLParamsRegex.FindAllString(path, -1)
		for i, p := range URLParams {
			URLParams[i] = p[1 : len(p)-1]
		}
	}
	return hasParams, isDelegate, URLParams
}
