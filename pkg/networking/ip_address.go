// Package networking contains various functions used to communicate between networks and
// draw data from the client network.
package networking

import (
	"net/http"
	"regexp"
	"strings"
)

// GetClientIP returns the current client's IP
func GetClientIP(req *http.Request) string {
	var addr string
	regex := regexp.MustCompile(`(?i)(?:for=)([^(;|,| )]+)`)

	if fwd := req.Header.Get(http.CanonicalHeaderKey("X-Forwarded-For")); fwd != "" {
		s := strings.Index(fwd, ", ")

		if s == -1 {
			s = len(fwd)
		}

		addr = fwd[:s]
	} else if fwd := req.Header.Get(http.CanonicalHeaderKey("X-Real-IP")); fwd != "" {
		addr = fwd
	} else if fwd := req.Header.Get(http.CanonicalHeaderKey("Forwarded")); fwd != "" {

		if match := regex.FindStringSubmatch(fwd); len(match) > 1 {
			addr = strings.Trim(match[1], `"`)
		}

	} else {
		addr = strings.Split(req.RemoteAddr, ":")[0]
	}

	return addr
}
