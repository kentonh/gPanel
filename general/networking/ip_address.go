package networking

import (
	"bytes"
	"net/http"
)

/*
	A function to get the client IP

	@return []byte
		The IP Address
	@return error
		If no error exists, the value is nil

	Known Problems: Uses an external API
*/
func GetClientIP() ([]byte, error) {
	resp, err := http.Get("http://myexternalip.com/raw")

	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.Bytes(), nil
}
