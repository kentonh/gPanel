// Package user is a child of package api to handle api calls concerning users
package user

import (
	"encoding/json"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/encryption"
)

// Register function is accessed by an API call from the webhost root
// by accessing /user_register and sending it a post request with userRequestData
// struct in JSON format.
func Register(res http.ResponseWriter, req *http.Request, dir string) bool {
	if req.Method != "POST" {
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	var userRequestData struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	err := json.NewDecoder(req.Body).Decode(&userRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	} else if len(userRequestData.User) == 0 || len(userRequestData.Pass) == 0 {
		http.Error(res, "Username or password field cannot be blank", http.StatusBadRequest)
		return false
	}

	ds, err := database.Open(dir + database.DB_USERS)
	if err != nil || ds == nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	var userDatabaseData struct {
		Pass   string `json:"pass"`
		Secret string `json:"secret"`
	}

	err = ds.Get(database.BUCKET_USERS, []byte(userRequestData.User), &userDatabaseData)
	if err != database.ErrKeyNotExist {
		http.Error(res, "Username already exists in the database", http.StatusBadRequest)
		return false
	}

	userDatabaseData.Pass, err = encryption.HashPassword(userRequestData.Pass)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	userDatabaseData.Secret = ""

	err = ds.Put(database.BUCKET_USERS, []byte(userRequestData.User), userDatabaseData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
