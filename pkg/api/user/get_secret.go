// Package user is a child of package api to handle api calls concerning users
package user

import "github.com/kentonh/gPanel/pkg/database"

// GetSecret is not accessible from the any client side request. It is
// only used on the server side to help verify users are who they say they
// are.
func GetSecret(user string, directory string) (string, error) {
	ds, err := database.Open(directory + database.DB_MAIN)
	if err != nil {
		return "", err
	}
	defer ds.Close()

	var userDatabaseData database.Struct_Users

	err = ds.Get(database.BUCKET_USERS, []byte(user), &userDatabaseData)
	if err != nil {
		return "", err
	}

	return userDatabaseData.Secret, nil
}
