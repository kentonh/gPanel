package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/kentonh/gPanel/pkg/database"
)

func List(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "GET" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, req.Method+" HTTP method is unsupported for this API.", http.StatusMethodNotAllowed)
		return false
	}

	ds, err := database.Open(dir + database.DB_MAIN)
	if err != nil || ds == nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	users, err := ds.ListAllUsers()
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	if len(users) > 0 {
		b, err := json.Marshal(users)
		if err != nil {
			logger.Println(req.URL.Path + "::" + err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return false
		}

		res.WriteHeader(http.StatusOK)
		res.Write(b)

		return true
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
