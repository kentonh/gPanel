package gpserver

import (
	"fmt"
	"io/ioutil"

	"github.com/Ennovar/gPanel/pkg/api/bundle"
	"github.com/Ennovar/gPanel/pkg/gpaccount"
)

func (con *Controller) detectBundles() {
	bundles := make(map[string]*gpaccount.Controller)

	dirs, err := ioutil.ReadDir("bundles/")
	if err != nil {
		fmt.Errorf("error finding bundles:%v", err.Error())
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			dirPath := "bundles/" + dir.Name() + "/"
			err, accPort, pubPort := bundle.GetPorts(dirPath)

			curBundle := gpaccount.New(dirPath, dir.Name(), accPort, pubPort)

			err = curBundle.Start()
			err2 := curBundle.Public.Start()
			if err != nil || err2 != nil {
				fmt.Println("error starting bundle:", dir.Name())
			}

			bundles[dir.Name()] = curBundle
		}
	}

	con.Bundles = bundles
}
