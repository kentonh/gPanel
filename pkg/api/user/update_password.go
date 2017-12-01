package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/encryption"
)

func UpdatePassword(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	var updatePasswordRequestData struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	err := json.NewDecoder(req.Body).Decode(&updatePasswordRequestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	ds, err := database.Open(dir + database.DB_MAIN)
	if err != nil || ds == nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	var userDatabaseData struct {
		Pass   string `json:"pass"`
		Secret string `json:"secret"`
	}

	userDatabaseData.Pass, err = encryption.HashPassword(updatePasswordRequestData.Pass)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	userDatabaseData.Secret = ""

	err = ds.Put(database.BUCKET_USERS, []byte(updatePasswordRequestData.User), userDatabaseData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
