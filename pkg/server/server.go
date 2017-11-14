package server

import (
	"io"
	"net/http"
	"os"

	"github.com/Ennovar/gPanel/pkg/account"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type Controller struct {
	Directory    string
	DocumentRoot string
	Bundles      []account.Controller
}

func New() *Controller {
	return &Controller{
		Directory:    "server/",
		DocumentRoot: "document_root/",
		Bundles:      []account.Controller{},
	}
}

func (con *Controller) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (con.Directory + con.DocumentRoot + "index.html")
	} else {
		path = (con.Directory + con.DocumentRoot + path)
	}

	if reqAuth(path) {
		if !con.checkAuth(res, req) {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}

	f, err := os.Open(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusNotFound, res)
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
		return
	}

	res.Header().Add("Content-Type", contentType)
	_, err = io.Copy(res, f)

	if err != nil {
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		return
	}
}
