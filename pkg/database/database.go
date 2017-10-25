// Package database handles all communication between software and the database
package database

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

func init() {
	db, err := bolt.Open("test.db", 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
