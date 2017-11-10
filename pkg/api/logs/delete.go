// Package logs is a child of package api to handle api calls concerning log files
package logs

import (
	"encoding/json"
	"net/http"
)

func Delete(res http.ResponseWriter, req *http.Request) bool {
	if req.Method != "GET" {
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

	return true
}
