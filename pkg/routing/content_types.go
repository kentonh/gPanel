// Package routing contains various functions related to HTTP routing
package routing

import (
	"errors"
	"path/filepath"
	"strings"
)

// getContentType returns the http header safe content-type attribute for the
// requested path. If the requested path is not formatted correctly with an extension
// then the function will return the zero-value for a string and an error.
// BUG(george-e-shaw-iv) Cannot handle interpreted files such as .php files
// BUG(george-e-shaw-iv) Does not cover the bulk of encountered content types on the web
func GetContentType(path string) (string, error) {
	var contentType string

	path = filepath.Base(path)
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)

	if ext == "" || fn == "" {
		return "", errors.New("Invalid file path")
	}

	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "text/javascript"
	case ".png":
		contentType = "image/png"
	case ".jpg":
		fallthrough
	case ".jpeg":
		contentType = "image/jpeg"
	default:
		contentType = "text/plain"
	}

	return contentType, nil
}
