package ssh

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/george-e-shaw-iv/nixtools"
)

func AddKey(res http.ResponseWriter, req *http.Request, logger *log.Logger) bool {
	if req.Method != "UPDATE" {
		logger.Println(req.URL.Path + "::" + req.Method + "::" + strconv.Itoa(http.StatusMethodNotAllowed) + "::" + http.StatusText(http.StatusMethodNotAllowed))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var requestData struct {
		Username  string `json:"username"`
		PublicKey string `json:"publickey"`
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

	if err = sysUser.AddAuthorizedKey(requestData.PublicKey); err != nil {
		logger.Println(req.URL.Path + "::" + err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusNoContent)
	return true
}
