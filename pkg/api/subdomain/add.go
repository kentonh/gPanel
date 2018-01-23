package subdomain

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
)

func Add(res http.ResponseWriter, req *http.Request, logger *log.Logger, dir string) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var requestData struct {
		Name string `json:"name"`
		Root string `json:"root"`
	}

	err := json.NewDecoder(req.Body).Decode(&requestData)
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

	err = ds.Put(database.BUCKET_SUBDOMAINS, []byte(requestData.Name), database.StructSubdomain{Root: requestData.Root})
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
