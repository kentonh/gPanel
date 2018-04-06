package nixtools

import (
	"os/exec"
	"bytes"
	"errors"
	"os"
	"io/ioutil"
	"strings"
)

// Function InitSSH creates the necessary folders,
// files, and generates a default key-pair for the
// given user. If parameter rootHasAccess is set
// to true then the public key of the root (sudo) user
// will be copied into the authorized_keys file of
// the user.
func (u *User) InitSSH(rootHasAccess bool) (error) {
	var stderr bytes.Buffer

	keygenArgs := []string{
		"-t",
		"rsa",
		"-N",
		"",
		"-f",
		"/home/" + u.Name + "/.ssh/id_rsa",
	}

	// Create the .ssh folder for said user
	cmd := exec.Command("mkdir", "/home/"+u.Name+"/.ssh")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	// Create authorized_keys file
	cmd = exec.Command("touch", "/home/"+u.Name+"/.ssh/authorized_keys")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	if rootHasAccess {
		// Put root public key into authorized_keys file
		cmd = exec.Command("cp", "/root/.ssh/id_rsa.pub", "/home/"+u.Name+"/.ssh/authorized_keys")
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return errors.New(stderr.String())
		}
	}

	// Create the default key-pair for said user
	cmd = exec.Command("ssh-keygen", keygenArgs...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	/* OWNERSHIP AND FILE PERMISSIONS START */
	cmd = exec.Command("chmod", "700", "/home/"+u.Name+"/.ssh")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	cmd = exec.Command("chmod", "600", "/home/"+u.Name+"/.ssh/id_rsa")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	cmd = exec.Command("chmod", "644", "/home/"+u.Name+"/.ssh/id_rsa.pub")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	cmd = exec.Command("chmod", "644", "/home/"+u.Name+"/.ssh/authorized_keys")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	cmd = exec.Command("chown", u.Name+":", "/home/"+u.Name+"/.ssh")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	cmd = exec.Command("chown", u.Name+":", "/home/"+u.Name+"/.ssh/id_rsa")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	cmd = exec.Command("chown", u.Name+":", "/home/"+u.Name+"/.ssh/id_rsa.pub")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}

	cmd = exec.Command("chown", u.Name+":", "/home/"+u.Name+"/.ssh/authorized_keys")
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}
	/* OWNERSHIP AND FILE PERMISSIONS END */

	return nil
}

// Function AddAuthorizedKey adds a new public key to
// a given user's authorized_keys file.
func (u *User) AddAuthorizedKey(key string) error {
	f, err := os.OpenFile("/home/"+u.Name+"/.ssh/authorized_keys", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err = f.WriteString(key + "\n"); err != nil {
		return err
	}

	return nil
}

// Function DeleteAuthorizedKey removes a public key
// that is already in the authorized_keys file of
// a given user.
func (u *User) DeleteAuthorizedKey(key string) error {
	old, err := ioutil.ReadFile("/home/" + u.Name + "/.ssh/authorized_keys")
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

	f, err := os.OpenFile("/home/"+u.Name+"/.ssh/authorized_keys", os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if _, err = f.WriteString(new); err != nil {
		return err
	}

	return nil
}

// Function GetAuthorizedKeys will return a slice
// of strings that contains all of the public keys
// within a given user's authorized_keys file.
// If the parameter removeRootKey is set to true the
// public key of the current root user of the system,
// if found within the file, will not be placed within
// the slice of strings.
func (u *User) GetAuthorizedKeys(removeRootKey bool) ([]string, error) {
	f, err := ioutil.ReadFile("/home/" + u.Name + "/.ssh/authorized_keys")
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

	// Remove root key if necessary
	if removeRootKey {
		if len(clean) == 1 {
			return nil, nil
		}

		key, err := ioutil.ReadFile("/root/.ssh/id_rsa.pub")
		if err != nil {
			return nil, err
		}

		for k, v := range clean {
			if v == string(key) {
				clean = append(clean[:k], clean[k+1:]...)
				break
			}
		}
	}

	return clean, nil
}