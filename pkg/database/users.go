package database

import "github.com/boltdb/bolt"

type Struct_Users struct {
	Pass   string `json:"pass"`
	Secret string `json:"secret"`
}

func (ds *Datastore) ListAllUsers() ([]string, error) {
	users := []string{}

	ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_USERS))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			users = append(users, string(k))
		}

		return nil
	})

	return users, nil
}
