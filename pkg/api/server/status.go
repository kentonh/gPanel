// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/public"
)

func Status(res http.ResponseWriter, req *http.Request, publicServer *public.Controller) bool {
	if req.Method != "GET" {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(strconv.Itoa(publicServer.Status)))

	return true
}
