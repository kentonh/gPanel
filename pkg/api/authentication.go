// Package api handles all API calls
package api

import (
	"encoding/json"
	"net/http"
)

// auth struct is the structure of the JSON data to be retrieved from
// the authentication API request
var auth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

// Authentication function is accessed by an API call from the webhost root
// by accessing /authentication and sending it a post request with
func Authentication(res http.ResponseWriter, req *http.Request) bool {
	if req.Method != "POST" {
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	err := json.NewDecoder(req.Body).Decode(&auth)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	} else {

		if auth.User == "root" && auth.Pass == "root" {
			res.WriteHeader(http.StatusNoContent)
			return true
		} else {
			http.Error(res, "Authentication failed", http.StatusUnauthorized)
			return false
		}

	}

}
