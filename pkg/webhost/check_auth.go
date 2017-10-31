// Package webhost handles the logic of the webhosting panel
package webhost

import "strings"

var allowedUnauthorizedPathSuffixes = [...]string{"api_testing.html", "user_auth", "user_register"}

func CheckAuth(path string) bool {
	for _, suffix := range allowedUnauthorizedPathSuffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	return false
}
