// Package logs is a child of package api to handle api calls concerning log files
package log

import (
	"encoding/json"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/file"
)

// Delete function is accessed from api/logs/delete and will attempt to
// delete a given log based off of request data.
func Delete(res http.ResponseWriter, req *http.Request, dir string) bool {
	if req.Method != "UPDATE" {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var deleteLogRequestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&deleteLogRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	var log string
	switch deleteLogRequestData.Name {
	case "client_errors":
		log = dir + "logs/" + file.LOG_CLIENT_ERRORS
	case "load_time":
		log = dir + "logs/" + file.LOG_LOADTIME
	case "server_errors":
		log = file.LOG_SERVER_ERRORS
	default:
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return false
	}

	handle, err := file.Open(log, true)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	closeErr, deleteErr := handle.Close(true)
	if closeErr != nil || deleteErr != nil {
		http.Error(res, closeErr.Error()+" AND "+deleteErr.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)

	return true
}
