// Package server is a child of package api to handle api calls concerning the server
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/kentonh/gPanel/pkg/public"
)

// Restart function is called from /api/server/restart and will attempt to shutdown, either gracefully
// or not gracefully (contingent on request data), and then turn back on the public server.
func Restart(res http.ResponseWriter, req *http.Request, logger *log.Logger, publicServer *public.Controller) bool {
	if req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var restartRequestData struct {
		graceful bool `json:"graceful"`
	}

	err := json.NewDecoder(req.Body).Decode(&restartRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	err = publicServer.Restart(restartRequestData.graceful)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	res.WriteHeader(http.StatusNoContent)

	return true
}
