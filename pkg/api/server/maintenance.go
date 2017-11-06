// Package server is a child of package api to handle api calls concerning the server
package server

import "net/http"

func Maintenance(res http.ResponseWriter, req *http.Request) bool {
	return true
}
