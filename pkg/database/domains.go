package database

import (
	"github.com/boltdb/bolt"
	"encoding/json"
)

type Struct_Domain struct {
	BundleName string `json:"name"`
	PublicPort int `json:"port"`
}

type Struct_Nameserver struct {
	Nameserver string `json:"nameserver"`
}

func (ds *Datastore) ListDomains(bundle string) (map[string]Struct_Domain, error) {
	filtered := make(map[string]Struct_Domain)
	var holder Struct_Domain

	ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_DOMAINS))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			json.Unmarshal(v, &holder)

			if bundle == "*" || holder.BundleName == bundle {
				filtered[string(k)] = holder
			}
		}

		return nil
	})

	return filtered, nil
}

func (ds *Datastore) RemoveInstances(bundle string) error {
	var holder Struct_Domain
	var toDelete []string

	ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_DOMAINS))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			json.Unmarshal(v, &holder)

			if holder.BundleName == bundle {
				toDelete = append(toDelete, string(k))
			}
		}

		return nil
	})

	for _, v := range toDelete {
		err := ds.Delete(BUCKET_DOMAINS, []byte(v))
		if err != nil && err != ErrKeyNotExist {
			return err
		}
	}

	return nil
}

func (ds *Datastore) ListNameservers() ([]string, error) {
	filtered := make([]string, 0)
	var holder Struct_Nameserver

	ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_NAMESERVERS))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			json.Unmarshal(v, &holder)

			filtered = append(filtered, holder.Nameserver)
		}

		return nil
	})

	return filtered, nil
}