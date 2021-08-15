package router

import (
	"errors"
	"strings"
)

func validatePath(path string) error {
	if !strings.HasSuffix(path, "/") {
		return errors.New("path must end with /")
	}
	if path[0] != '/' {
		return errors.New("path must start with /")
	}
	if strings.Contains(path, "//") {
		return errors.New("path must not contain //")
	}
	if hasSpaceRegex.FindString(path) != "" {
		return errors.New("path must not contain spaces")
	}
	if validatePathRegex.MatchString(path) {
		return nil
	}
	return errors.New("invalid path. path should starts and ends with  / ")
}
