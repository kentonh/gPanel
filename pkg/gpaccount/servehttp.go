package gpaccount

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/routing"
)

// ServeHTTP function routes all requests for the private webhost server. It is used in the main
// function inside of the http.ListenAndServe() function for the private webhost host.
func (con *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (con.DocumentRoot + "index.html")
	} else {
		path = (con.DocumentRoot + path)
	}

	if reqAuth(path) {
		if !con.checkAuth(res, req) {
			con.AccountLogger.Println(path + "::" + strconv.Itoa(http.StatusUnauthorized) + "::" + http.StatusText(http.StatusUnauthorized))
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
		con.AccountLogger.Println(path + "::" + strconv.Itoa(http.StatusNotFound) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusNotFound, res)
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		con.AccountLogger.Println(path + "::" + strconv.Itoa(http.StatusUnsupportedMediaType) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
		return
	}

	res.Header().Add("Content-Type", contentType)
	_, err = io.Copy(res, f)

	if err != nil {
		con.AccountLogger.Println(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		return
	}
}
