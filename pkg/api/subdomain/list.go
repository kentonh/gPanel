package subdomain

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
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	ds, err := database.Open(dir + database.DB_MAIN)
	if err != nil || ds == nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	subdomains, err := ds.ListSubdomains()
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	if len(subdomains) > 0 {
		b, err := json.Marshal(subdomains)
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
