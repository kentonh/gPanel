package database

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type StructSubdomain struct {
	Root string `json:"root"`
}

func (ds *Datastore) ListSubdomains() (map[string]StructSubdomain, error) {
	filtered := make(map[string]StructSubdomain)
	var holder StructSubdomain

	ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_SUBDOMAINS))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			json.Unmarshal(v, &holder)
			filtered[string(k)] = holder
		}

		return nil
	})

	return filtered, nil
}
