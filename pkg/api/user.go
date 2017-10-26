// Package api handles all API calls
package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/database"
)

// userRequestData struct is the structure of the JSON data to be
// retrieved from the authentication API request
var userRequestData struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

// userDatabaseData struct is the structure of the JSON data to be retrieved from
// the bolt database inside of the user bucket, using the username as the key.
var userDatabaseData struct {
	Pass string `json:"pass"`
}

// UserAuthentication function is accessed by an API call from the webhost root
// by accessing /authentication and sending it a post request with
func UserAuthentication(res http.ResponseWriter, req *http.Request) bool {
	if req.Method != "POST" {
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	err := json.NewDecoder(req.Body).Decode(&userRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	ds, err := database.Open(database.DBLOC_MAIN)
	if err != nil || ds == nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}
	defer ds.Close()

	err = ds.Get(database.BUCKET_USERS, []byte(userRequestData.User), &userDatabaseData)

	if err == database.ErrKeyNotExist {
		http.Error(res, "User does not exist in database", http.StatusUnauthorized)
		return false
	}

	if userRequestData.Pass != userDatabaseData.Pass {
		http.Error(res, "Invalid password", http.StatusUnauthorized)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true

}

// UserAuthentication function is accessed by an API call from the webhost root
// by accessing /authentication and sending it a post request with
func UserRegistration(res http.ResponseWriter, req *http.Request) bool {
	if req.Method != "POST" {
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	err := json.NewDecoder(req.Body).Decode(&userRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	} else if len(userRequestData.User) == 0 || len(userRequestData.Pass) == 0 {
		http.Error(res, "Username or password field cannot be blank", http.StatusBadRequest)
		return false
	}

	ds, err := database.Open(database.DBLOC_MAIN)
	if err != nil || ds == nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}
	defer ds.Close()

	userDatabaseData.Pass = userRequestData.Pass
	err = ds.Put(database.BUCKET_USERS, []byte(userRequestData.User), userDatabaseData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true

}
