package router

import (
	"fmt"
	"strings"
)

func (path Path) String() string {
	return string(path)
}
func (path *Path) Validate() {
	*path = Path(strings.TrimSpace(path.String()))
	// should not contain  // or /../
	if strings.Contains(path.String(), "//") || strings.Contains(path.String(), ".") {
		panic("path must not include // or .")
	}
	if path.String() == "" {
		panic("path must not be empty")
	}
	if !(strings.HasPrefix(path.String(), "/") && (strings.HasSuffix(path.String(), "/") || path.String() == "/")) {
		panic(fmt.Sprintf("path %s must start with / and end with /", path.String()))
	}
}

// // isParamKey checks if param key is duplicated
// func isParamKey(params []string, key string) bool {
// 	for _, v := range params {
// 		if len(key) <= 1 {
// 			return false
// 		}
// 		if v == key[1:] {
// 			return true
// 		}
// 	}
// 	return false
// }

// func prepareRequestPath(path string) string {
// 	if path == "" {
// 		path = "/"
// 	}
// 	if path != "/" && len(path) > 1 {
// 		if !validateRequestPathRegex.MatchString(path) {
// 			path = fmt.Sprintf("%s/", path)
// 		}
// 	}
// 	return path
// }

// func getPathInfo(path string) (hasParams, isDelegate bool, URLParams []string) {
// 	isDelegate = delegateRegex.MatchString(path)
// 	hasParams = hasParamsRegex.MatchString(path)
// 	if hasParams {
// 		URLParams = getURLParamsRegex.FindAllString(path, -1)
// 		for i, p := range URLParams {
// 			URLParams[i] = p[1 : len(p)-1]
// 		}
// 	}
// 	return hasParams, isDelegate, URLParams
// }
