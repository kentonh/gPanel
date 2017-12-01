package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
)

func Delete(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	var deleteUserRequestData struct {
		User string `json:"user"`
	}

	err := json.NewDecoder(req.Body).Decode(&deleteUserRequestData)
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

	count, err := ds.Count(database.BUCKET_USERS)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	if count <= 1 {
		logger.Println(req.URL.Path + ":: if only one user exists it cannot be deleted")
		http.Error(res, "If only one user exists it cannot be deleted", http.StatusBadRequest)
		return false
	}

	err = ds.Delete(database.BUCKET_USERS, []byte(deleteUserRequestData.User))
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
