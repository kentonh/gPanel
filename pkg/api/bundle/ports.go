package bundle

import (
	"os"

	"github.com/kentonh/gPanel/pkg/database"
)

func GetPorts(dir string) (error, int, int) {
	if _, err := os.Stat(dir + "datastore.db"); os.IsNotExist(err) {
		return err, 0, 0
	}

	ds, err := database.Open(dir + database.DB_MAIN)
	if err != nil {
		return err, 0, 0
	}
	defer ds.Close()

	var databaseBundlePorts struct {
		Account int `json:"account"`
		Public  int `json:"public"`
	}

	err = ds.Get(database.BUCKET_PORTS, []byte("bundle_ports"), &databaseBundlePorts)
	if err != nil {
		return err, 0, 0
	}

	return nil, databaseBundlePorts.Account, databaseBundlePorts.Public
}
