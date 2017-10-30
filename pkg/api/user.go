// Package api handles all API calls
package api

import (
	"encoding/json"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/encryption"
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
// by accessing /user_auth and sending it a post request with userRequestData
// struct in JSON format.
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

	err = encryption.CheckPassword([]byte(userDatabaseData.Pass), []byte(userRequestData.Pass))
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}

// UserRegistration function is accessed by an API call from the webhost root
// by accessing /user_register and sending it a post request with userRequestData
// struct in JSON format.
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

	err = ds.Get(database.BUCKET_USERS, []byte(userRequestData.User), &userDatabaseData)
	if err != database.ErrKeyNotExist {
		http.Error(res, "Username already exists in the database", http.StatusBadRequest)
		return false
	}

	userDatabaseData.Pass, err = encryption.HashPassword(userRequestData.Pass)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	err = ds.Put(database.BUCKET_USERS, []byte(userRequestData.User), userDatabaseData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
