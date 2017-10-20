// Package networking contains various functions used to communicate between networks and
// draw data from the client network.
package networking

import (
	"bytes"
	"net/http"
)

// GetClientIP returns the current client's IP as an array of bytes.
// BUG(george-e-shaw-iv) Uses an external API
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
