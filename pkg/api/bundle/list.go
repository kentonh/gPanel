package bundle

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/kentonh/gPanel/pkg/gpaccount"
)

func List(res http.ResponseWriter, req *http.Request, logger *log.Logger, bundles map[string]*gpaccount.Controller) bool {
	if req.Method != "GET" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusNotFound) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	if len(bundles) <= 0 {
		res.WriteHeader(http.StatusNoContent)
		return true
	}

	keys := make([]string, 0, len(bundles))
	for key := range bundles {
		keys = append(keys, key)
	}

	jsonData, err := json.Marshal(keys)
	if err != nil {
		logger.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

	return true
}
