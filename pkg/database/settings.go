package database

import (
	"encoding/json"
	"errors"

	"github.com/kentonh/gPanel/pkg/emailer"
	"github.com/boltdb/bolt"
)

type Struct_SMTP struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
}

type Struct_Admin struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Function CheckAdminSettings makes sure that the
// admin settings are set and are valid.
func (ds *Datastore) CheckAdminSettings() (rerr error) {
	rerr = ds.handle.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BUCKET_GENERAL))
		smtp := b.Get([]byte("smtp"))

		var smtpCreds Struct_SMTP
		err := json.Unmarshal(smtp, &smtpCreds)
		if err != nil {
			return err
		}

		_, err = emailer.New(smtpCreds.Type, emailer.Credentials{
			Username: smtpCreds.Username,
			Password: smtpCreds.Password,
			Server:   smtpCreds.Server,
			Port:     smtpCreds.Port,
		})

		if err != nil {
			return err
		}

		a := b.Get([]byte("admin"))

		var admin Struct_Admin
		err = json.Unmarshal(a, &admin)
		if err != nil {
			return err
		}

		if len(admin.Email) == 0 || len(admin.Name) == 0 {
			return errors.New("admin name and email settings are empty")
		}

		return nil
	})

	return
}
