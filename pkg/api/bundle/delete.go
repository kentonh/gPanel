package bundle

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/gpaccount"
	"github.com/Ennovar/gPanel/pkg/system"
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

	ds, err := database.Open("server/" + database.DB_DOMAINS)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}
	defer ds.Close()

	err = ds.RemoveInstances(rqData.Name)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
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

	err, err2 := system.DeleteBundleUser(rqData.Name)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error() + " AND " + err2.Error())
		http.Error(res, err2.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
