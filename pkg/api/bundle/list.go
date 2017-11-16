package bundle

import (
	"encoding/json"
	"net/http"

	"github.com/Ennovar/gPanel/pkg/gpaccount"
)

func List(res http.ResponseWriter, req *http.Request, bundles map[string]*gpaccount.Controller) bool {
	if req.Method != "GET" {
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
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

	return true
}
