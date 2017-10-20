// Package routing contains various functions related to HTTP routing
package routing

import "net/http"

func HttpThrowStatus(code int, context http.ResponseWriter) {
	context.WriteHeader(code)
	context.Write([]byte(http.StatusText(code)))
}
