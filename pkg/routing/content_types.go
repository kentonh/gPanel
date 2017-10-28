// Package routing contains various functions related to HTTP routing
package routing

import (
	"errors"
	"mime"
	"path/filepath"
	"strings"
)

// getContentType returns the http header safe Content-Type for the requested
// path. If the requested path is not formatted correctly with an extension
// then the function will return the zero-value for a string and an error.
func GetContentType(path string) (string, error) {
	var ct string

	path = filepath.Base(path)
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)

	if ext == "" || fn == "" {
		return "", errors.New("Invalid file path")
	}

	ct = mime.TypeByExtension(ext)

	if ct == "" {
		return "", errors.New("Could not get content type")
	}

	return ct, nil
}
