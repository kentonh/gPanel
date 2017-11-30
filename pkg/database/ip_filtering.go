package database

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
)

func IDtoKey(id int) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(id))
	return key
}

func (ds *Datastore) NewFilteredIP(newip *Struct_Filtered_IP) error {
	return ds.handle.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_FILTERED_IPS))

		id, _ := b.NextSequence()
		newip.ID = int(id)

		buf, err := json.Marshal(newip)
		if err != nil {
			return err
		}

		key := IDtoKey(newip.ID)
		return b.Put(key, buf)
	})
}

func (ds *Datastore) GetFilteredIPs(iptype string) (map[int]Struct_Filtered_IP, error) {
	filtered := make(map[int]Struct_Filtered_IP)
	var holder Struct_Filtered_IP

	ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_FILTERED_IPS))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			json.Unmarshal(v, &holder)

			if holder.Type == iptype {
				filtered[holder.ID] = holder
			}
		}

		return nil
	})

	return filtered, nil
}

func (ds *Datastore) IsFiltered(ip string, iptype string) (bool, error) {
	var holder Struct_Filtered_IP
	found := false

	ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_FILTERED_IPS))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			json.Unmarshal(v, &holder)

			if holder.IP == ip && holder.Type == iptype {
				found = true
				break
			}
		}

		return nil
	})

	return found, nil
}
