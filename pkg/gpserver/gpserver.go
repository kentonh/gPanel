// Package gpserver handles the logic of the gPanel server
package gpserver

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api/bundle"
	"github.com/Ennovar/gPanel/pkg/gpaccount"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type Controller struct {
	Directory    string
	DocumentRoot string
	Bundles      map[string]*gpaccount.Controller
}

func New() *Controller {
	bundles := make(map[string]*gpaccount.Controller)

	dirs, err := ioutil.ReadDir("bundles/")
	if err != nil {
		log.Fatal("Error finding bundles: %v\n", err.Error())
	}

	for _, dir := range dirs {
		if dir.Name() == "default_bundle" || !dir.IsDir() {
			continue
		}

		if strings.HasPrefix(dir.Name(), "bundle_") {
			dirPath := "bundles/" + dir.Name() + "/"
			err, accPort, pubPort := bundle.GetPorts(dirPath)

			curBundle := gpaccount.New(dirPath, accPort, pubPort)

			err = curBundle.Start()
			err2 := curBundle.Public.Start()
			if err != nil || err2 != nil {
				log.Fatal("Error starting bundle: %v\n", dir.Name())
			}

			bundles[strings.Replace(dir.Name(), "bundle_", "", 1)] = curBundle
		}
	}

	return &Controller{
		Directory:    "server/",
		DocumentRoot: "document_root/",
		Bundles:      bundles,
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

	isApi, _ := con.apiHandler(res, req, 0)

	if isApi {
		// API methods handle HTTP logic from here
		return
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
