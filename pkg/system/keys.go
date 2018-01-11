package system

import (
	"io/ioutil"
	"os"
	"strings"
)

func AddAuthorizedKey(username, key string) error {
	f, err := os.OpenFile("/home/"+username+"/.ssh/authorized_keys", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err = f.WriteString(key); err != nil {
		return err
	}

	return nil
}

func DeleteAuthorizedKey(username, key string) error {
	old, err := ioutil.ReadFile("/home/" + username + "/.ssh/authorized_keys")
	if err != nil {
		return err
	}

	lines := strings.Split(string(old), "\n")

	for i, line := range lines {
		if strings.Contains(line, key) {
			lines = append(lines[:i], lines[i+1:]...)
			break
		}
	}

	new := strings.Join(lines, "\n")

	f, err := os.OpenFile("/home/"+username+"/.ssh/authorized_keys", os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err = f.WriteString(new); err != nil {
		return err
	}

	return nil
}

func GetAuthorizedKeys(username string) ([]string, error) {
	f, err := ioutil.ReadFile("/home/" + username + "/.ssh/authorized_keys")
	if err != nil {
		return nil, err
	}

	// Remove empty strings
	raw := strings.Split(string(f), "\n")
	var clean []string

	for _, each := range raw {
		if each != "" {
			clean = append(clean, each)
		}
	}

	// Remove root key
	if len(clean) == 1 {
		return nil, nil
	}

	return clean[1:], nil
}
