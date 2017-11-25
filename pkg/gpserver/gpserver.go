// Package gpserver handles the logic of the gPanel server
package gpserver

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api/bundle"
	"github.com/Ennovar/gPanel/pkg/file"
	"github.com/Ennovar/gPanel/pkg/gpaccount"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type Controller struct {
	Directory    string
	DocumentRoot string
	Bundles      map[string]*gpaccount.Controller
	ServerLogger *file.Handler
}

func New() *Controller {
	bundles := make(map[string]*gpaccount.Controller)

	dirs, err := ioutil.ReadDir("bundles/")
	if err != nil {
		fmt.Errorf("Error finding bundles: %v\n", err.Error())
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
				fmt.Errorf("Error starting bundle: %v\n", dir.Name())
			}

			bundles[strings.Replace(dir.Name(), "bundle_", "", 1)] = curBundle
		}
	}

	serverErrorLogger, err := file.Open(file.LOG_SERVER_ERRORS, true)
	if err != nil {
		fmt.Errorf("Error whilst trying to start server logging instance: %v\n", err.Error())
	}

	return &Controller{
		Directory:    "server/",
		DocumentRoot: "document_root/",
		Bundles:      bundles,
		ServerLogger: serverErrorLogger,
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
			con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusUnauthorized) + "::" + http.StatusText(http.StatusUnauthorized))
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
		con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusNotFound) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusNotFound, res)
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusUnsupportedMediaType) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
		return
	}

	res.Header().Add("Content-Type", contentType)
	_, err = io.Copy(res, f)

	if err != nil {
		con.ServerLogger.Write(path + "::" + strconv.Itoa(http.StatusInternalServerError) + "::" + err.Error())
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		return
	}
}
