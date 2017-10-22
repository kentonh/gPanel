// Package routing contains various functions related to HTTP routing
package routing

import (
	"errors"
	"path/filepath"
)

// getContentType returns the http header safe content-type attribute for the
// requested path. If the requested path is not formatted correctly with an extension
// then the function will return the zero-value for a string and an error.
// BUG(george-e-shaw-iv) Cannot handle interpreted files such as .php files
// BUG(george-e-shaw-iv) Does not cover the bulk of encountered content types on the web
func GetContentType(path string) (string, error) {
	var contentType string
	fileType := filepath.Ext(path)

	if fileType == "" {
		return contentType, errors.New("Invalid path, contained no period-separated content-type")
	}

	switch fileType {
	case "html":
		contentType = "text/html"
	case "css":
		contentType = "text/css"
	case "js":
		contentType = "text/javascript"
	case "png":
		contentType = "image/png"
	case "jpg":
		fallthrough
	case "jpeg":
		contentType = "image/jpeg"
	default:
		contentType = "text/plain"
	}

	return contentType, nil
}
