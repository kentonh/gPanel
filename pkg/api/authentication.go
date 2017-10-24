// Package api handles all API calls
package api

import (
	"encoding/json"
	"net/http"
)

var auth struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func Authentication(res http.ResponseWriter, req *http.Request) bool {
	err := json.NewDecoder(req.Body).Decode(&auth)

	if err != nil {
		http.Error(res, err.Error(), 400)
		return false
	} else {

		if auth.User == "root" && auth.Pass == "root" {
			res.WriteHeader(200)
			res.Write([]byte("success"))
			return true
		} else {
			http.Error(res, "Authentication failed", 401)
			return false
		}

	}

}
