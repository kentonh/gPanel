package bundle

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/file"
	"github.com/Ennovar/gPanel/pkg/gpaccount"
)

func Create(res http.ResponseWriter, req *http.Request, bundles map[string]*gpaccount.Controller) bool {
	if req.Method != "POST" {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}

	var createBundleRequestData struct {
		Name    string `json:"name"`
		AccPort int    `json:"account_port"`
		PubPort int    `json:"public_port"`
	}

	err := json.NewDecoder(req.Body).Decode(&createBundleRequestData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	check, err := net.Listen("tcp", ":"+strconv.Itoa(createBundleRequestData.AccPort))
	if err != nil {
		http.Error(res, "A service is already listening on port "+strconv.Itoa(createBundleRequestData.AccPort), http.StatusInternalServerError)
		return false
	}
	check.Close()

	check, err = net.Listen("tcp", ":"+strconv.Itoa(createBundleRequestData.PubPort))
	if err != nil {
		http.Error(res, "A service is already listening on port "+strconv.Itoa(createBundleRequestData.PubPort), http.StatusInternalServerError)
		return false
	}
	check.Close()

	err = nil
	for k, v := range bundles {
		if k == createBundleRequestData.Name {
			err = errors.New("Bundle \"" + k + "\" already exists")
			break
		}

		if v.Port == createBundleRequestData.AccPort ||
			v.Port == createBundleRequestData.PubPort ||
			v.Public.Port == createBundleRequestData.AccPort ||
			v.Public.Port == createBundleRequestData.PubPort {
			err = errors.New("An existing bundle is using the port \"" + strconv.Itoa(v.Port) + "\" already")
			break
		}
	}

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	newBundle := "bundles/bundle_" + createBundleRequestData.Name
	err = file.CopyDir("bundles/default_bundle", newBundle)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false
	}

	ds, err := database.Open(newBundle + "/" + database.DB_MAIN)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}
	defer ds.Close()

	err = ds.Put(database.BUCKET_PORTS, []byte("account"), createBundleRequestData.AccPort)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	err = ds.Put(database.BUCKET_PORTS, []byte("public"), createBundleRequestData.PubPort)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return false
	}

	bundles[createBundleRequestData.Name] = gpaccount.New(newBundle+"/", createBundleRequestData.AccPort, createBundleRequestData.PubPort)
	_ = bundles[createBundleRequestData.Name].Start()
	_ = bundles[createBundleRequestData.Name].Public.Start()

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(createBundleRequestData.Name))

	return true
}
