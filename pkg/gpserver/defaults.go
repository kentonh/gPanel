package gpserver

import (
	"fmt"

	"github.com/Ennovar/gPanel/pkg/database"
	"github.com/Ennovar/gPanel/pkg/encryption"
)

func (con *Controller) setDefaults() {
	ds, err := database.Open(con.Directory + database.DB_MAIN)
	if err != nil || ds == nil {
		fmt.Println("error whilst trying to set server defaults:", err.Error())
		return
	}
	defer ds.Close()

	users, err := ds.Count(database.BUCKET_USERS)
	if users >= 1 {
		return
	}

	var defaults database.Struct_Users

	defaults.Pass, err = encryption.HashPassword("root")
	if err != nil {
		fmt.Println("error whilst trying to set server defaults:", err.Error())
		return
	}
	defaults.Secret = ""

	err = ds.Put(database.BUCKET_USERS, []byte("root"), defaults)
	if err != nil {
		fmt.Println("error whilst trying to set server defaults:", err.Error())
		return
	}

	fmt.Print("Since there are no stored users for the server upon startup the default user \"root\" has been set with the password \"root\"\n\n")
	fmt.Print("Upon your first time logging into the gPanel Server please either create a new user and delete the user root, or change the user root's password to something more secure.\n\n")
}
