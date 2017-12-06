package bundle

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/gpaccount"
)

func Delete(res http.ResponseWriter, req *http.Request, logger *log.Logger, bundles map[string]*gpaccount.Controller, dir string) bool {
	if req.Method != "DELETE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var rqData struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&rqData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	if _, ok := bundles[rqData.Name]; !ok {
		logger.Println(req.URL.Path + ":: bundle does not exist")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	_ = bundles[rqData.Name].Public.Stop(false)
	_ = bundles[rqData.Name].Stop(false)

	err = os.RemoveAll(dir)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	delete(bundles, rqData.Name)

	res.WriteHeader(http.StatusNoContent)
	return true
}
