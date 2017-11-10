// Package logs is a child of package api to handle api calls concerning log files
package logs

import (
	"encoding/json"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/file"
)

func Read(res http.ResponseWriter, req *http.Request) bool {
	if req.Method != "POST" {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var readLogRequestData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&readLogRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	valid := false
	for _, log := range file.KNOWN_LOGS {
		if log == readLogRequestData.Name {
			valid = true
			break
		}
	}

	if !valid {
		http.Error(res, "Attempting to read from an unknown log file.", http.StatusBadRequest)
		return false
	}

	handle, err := file.Open(readLogRequestData.Name, true, true)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}
	defer handle.Close(false)

	data, err := handle.Read()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Write(data)

	return true
}
