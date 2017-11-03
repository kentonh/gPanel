// Package user is a child of package api to handle api calls concerning users
package user

// userRequestData struct is the structure of the JSON data to be
// retrieved from the authentication API request
var userRequestData struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

// userDatabaseData struct is the structure of the JSON data to be retrieved from
// the bolt database inside of the user bucket, using the username as the key.
var userDatabaseData struct {
	Pass   string `json:"pass"`
	Secret string `json:"secret"`
}
