// Package user is a child of package api to handle api calls concerning users
package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/kentonh/gPanel/pkg/database"
	"github.com/kentonh/gPanel/pkg/encryption"
)

// Register function is accessed by an API call from the webhost root
// by accessing /user_register and sending it a post request with userRequestData
// struct in JSON format.
func Register(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	var userRequestData struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	err := json.NewDecoder(req.Body).Decode(&userRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	} else if len(userRequestData.Pass) < 8  {
		logger.Println(req.URL.Path + "::password must be at least 8 characters long")
		http.Error(res, "Password must be at least 8 characters long!", http.StatusBadRequest)
		return false
	} else if len(userRequestData.User) == 0 || len(userRequestData.Pass) == 0 {
		logger.Println(req.URL.Path + "::username or password field cannot be blank")
		http.Error(res, "Username or password field cannot be blank", http.StatusBadRequest)
		return false
	}

	ds, err := database.Open(dir + database.DB_MAIN)
	if err != nil || ds == nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	var userDatabaseData database.Struct_Users

	err = ds.Get(database.BUCKET_USERS, []byte(userRequestData.User), &userDatabaseData)
	if err != database.ErrKeyNotExist {
		logger.Println(req.URL.Path + "::username already exists in the database")
		http.Error(res, "Username already exists in the database", http.StatusBadRequest)
		return false
	}

	userDatabaseData.Pass, err = encryption.HashPassword(userRequestData.Pass)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	userDatabaseData.Secret = ""

	err = ds.Put(database.BUCKET_USERS, []byte(userRequestData.User), userDatabaseData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
