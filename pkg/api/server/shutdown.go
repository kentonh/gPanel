// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"encoding/json"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/public"
)

// Shutdown function is called from /api/server/shutdown and will attempt to shutdown, either gracefully
// or not gracefully (contingent on request data).
func Shutdown(res http.ResponseWriter, req *http.Request, publicServer *public.Controller) bool {
	if req.Method != "UPDATE" {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var shutdownRequestData struct {
		Graceful bool `json:"graceful"`
	}

	err := json.NewDecoder(req.Body).Decode(&shutdownRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	err = publicServer.Stop(shutdownRequestData.Graceful)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	res.WriteHeader(http.StatusNoContent)

	return true
}
