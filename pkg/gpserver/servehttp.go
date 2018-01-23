package gpserver

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"strings"

	"github.com/Ennovar/gPanel/pkg/routing"
)

func (con *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (con.Directory + con.DocumentRoot + "index.html")
	} else {
		path = (con.Directory + con.DocumentRoot + path)
	}

	if strings.HasSuffix(path, "throw400") {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if strings.HasSuffix(path, "index.html") {
		if con.checkAuth(res, req) == true {
			http.Redirect(res, req, "gPanel.html", http.StatusFound)
		}
	}

	if reqAuth(path) {
		if !con.checkAuth(res, req) {
			con.ServerLogger.Println(path + "::" + strconv.Itoa(http.StatusUnauthorized) + "::" + http.StatusText(http.StatusUnauthorized))
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}

	isApi, _ := con.apiHandler(res, req)

	if isApi {
		// API methods handle HTTP logic from here
		return
	}

	f, err := os.Open(path)

	if err != nil {
		con.ServerLogger.Println(path + "::" + strconv.Itoa(http.StatusNotFound) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusNotFound, res)
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		con.ServerLogger.Println(path + "::" + strconv.Itoa(http.StatusUnsupportedMediaType) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
		return
	}

	res.Header().Add("Content-Type", contentType)
	_, err = io.Copy(res, f)

	if err != nil {
		con.ServerLogger.Println(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		return
	}
}
