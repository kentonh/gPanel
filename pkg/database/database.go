// Package database handles all communication between software and the database
package database

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

// Database constants
const (
	DB_MAIN     = "datastore.db"
	DB_SETTINGS = "settings.db"
)

// Bucket constants
const (
	// DB_MAIN BUCKETS
	BUCKET_USERS        = "users"
	BUCKET_PORTS        = "ports"
	BUCKET_FILTERED_IPS = "filtered_ips"

	// DB_SETTINGS BUCKETS
	BUCKET_GENERAL = "general"
)

// Error codes
var (
	ErrKeyNotExist = errors.New("key does not exist")
)

type Datastore struct {
	handle *bolt.DB
}

// Open function will open the database and return a Datastore struct
// that has a handle within it for various datastore functions.
func Open(filepath string) (*Datastore, error) {
	db, err := bolt.Open(filepath, 0666, &bolt.Options{Timeout: 15 * time.Second})

	if err != nil {
		return nil, err
	}

	ds := &Datastore{
		handle: db,
	}

	// Ensure that all top-level buckets exist
	err = ds.handle.Update(func(tx *bolt.Tx) error {
		if strings.HasSuffix(filepath, DB_MAIN) {
			_, err := tx.CreateBucketIfNotExists([]byte(BUCKET_USERS))
			if err != nil {
				return err
			}

			_, err = tx.CreateBucketIfNotExists([]byte(BUCKET_PORTS))
			if err != nil {
				return err
			}

			_, err = tx.CreateBucketIfNotExists([]byte(BUCKET_FILTERED_IPS))
			if err != nil {
				return err
			}
		}

		if strings.HasSuffix(filepath, DB_SETTINGS) {
			_, err = tx.CreateBucketIfNotExists([]byte(BUCKET_GENERAL))
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ds, nil
}

// Close is a function attached to the Datastore struct that will close the current bolt instance
func (ds *Datastore) Close() error {
	return ds.handle.Close()
}

// Get is a function attached to the Datastore struct that will get data from the current bolt instance
func (ds *Datastore) Get(bucket string, key []byte, result interface{}) error {
	return ds.handle.View(func(tx *bolt.Tx) error {
		dsValue := tx.Bucket([]byte(bucket)).Get(key)

		if dsValue == nil {
			return ErrKeyNotExist
		}

		return json.Unmarshal(dsValue, result)
	})
}

// Put is a function attached to the Datastore struct that will put data into the current bolt instance
func (ds *Datastore) Put(bucket string, key []byte, value interface{}) error {
	var err error
	dsValue, ok := value.([]byte)

	if !ok {
		dsValue, err = json.Marshal(value)

		if err != nil {
			return err
		}

	}

	return ds.handle.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Put(key, dsValue)
	})
}

// Delete is a function attached to the Datastore struct that will delete data from the current bolt instance
func (ds *Datastore) Delete(bucket string, key []byte) error {
	return ds.handle.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Delete(key)
	})
}

func (ds *Datastore) Count(bucket string) (int, error) {
	count := 0

	ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			count++
		}

		return nil
	})

	return count, nil
}
