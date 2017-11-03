// Package user is a child of package api to handle api calls concerning users
package user

import "github.com/Ennovar/gPanel/pkg/database"

// UserSecret is not accessible from the any client side request. It is
// only used on the server side to help verify users are who they say they
// are.
func UserSecret(user string) (string, error) {
	ds, err := database.Open(database.DBLOC_MAIN)
	if err != nil {
		return "", err
	}

	err = ds.Get(database.BUCKET_USERS, []byte(user), &userDatabaseData)
	if err != nil {
		return "", err
	}

	return userDatabaseData.Secret, nil
}
