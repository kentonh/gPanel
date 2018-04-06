package ssh

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/george-e-shaw-iv/nixtools"
)

func GetKeys(res http.ResponseWriter, req *http.Request, logger *log.Logger) bool {
	if req.Method != "POST" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var requestData struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(req.Body).Decode(&requestData)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	sysUser, err := nixtools.GetUser(requestData.Username, false)
	if err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	var keys []string
	if keys, err = sysUser.GetAuthorizedKeys(true); err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	if len(keys) > 0 {
		b, err := json.Marshal(keys)
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
