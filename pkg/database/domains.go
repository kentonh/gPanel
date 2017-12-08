package database

import (
	"github.com/boltdb/bolt"
	"encoding/json"
)

type Struct_Domain struct {
	BundleName string `json:"name"`
}

func (ds *Datastore) ListDomains(bundle string) (map[string]string, error) {
	filtered := make(map[string]string)
	var holder Struct_Domain

	ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_DOMAINS))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			json.Unmarshal(v, &holder)

			if bundle == "*" || holder.BundleName == bundle {
				filtered[holder.BundleName] = string(k)
			}
		}

		return nil
	})

	return filtered, nil
}