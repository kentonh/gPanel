package system

import (
	"bytes"
	"errors"
	"os/exec"
)

func CreateBundleUser(username string) (error, error) {
	var cerr bytes.Buffer
	var err error

	adduserArgs := []string{
		"--disabled-password",
		"--gecos",
		"",
		username,
	}

	keygenArgs := []string{
		"-t",
		"rsa",
		"-N",
		"",
		"-f",
		"/home/" + username + "/.ssh/id_rsa",
	}

	// Add the user
	cmd := exec.Command("adduser", adduserArgs...)
	cmd.Stderr = &cerr

	if err = cmd.Run(); err != nil {
		return err, errors.New(cerr.String())
	}

	// Create the .ssh folder for said user
	cmd = exec.Command("mkdir", "/home/"+username+"/.ssh")
	cmd.Stderr = &cerr

	if err = cmd.Run(); err != nil {
		return err, errors.New(cerr.String())
	}

	// Create authorized_keys file
	cmd = exec.Command("touch", "/home/"+username+"/.ssh/authorized_keys")
	cmd.Stderr = &cerr

	if err = cmd.Run(); err != nil {
		return err, errors.New(cerr.String())
	}

	// Put root public key into authorized_keys file
	cmd = exec.Command("cp", "/root/.ssh/id_rsa.pub", "/home/"+username+"/.ssh/authorized_keys")
	cmd.Stderr = &cerr

	if err = cmd.Run(); err != nil {
		return err, errors.New(cerr.String())
	}

	// Correct permissions on authorized_keys file
	cmd = exec.Command("chmod", "600", "/home/"+username+"/.ssh/authorized_keys")
	cmd.Stderr = &cerr

	if err = cmd.Run(); err != nil {
		return err, errors.New(cerr.String())
	}

	// Create the host key-pair for said user
	cmd = exec.Command("ssh-keygen", keygenArgs...)
	cmd.Stderr = &cerr

	if err = cmd.Run(); err != nil {
		return err, errors.New(cerr.String())
	}

	// Create document root for said user
	cmd = exec.Command("mkdir", "/home/"+username+"/document_root")
	cmd.Stderr = &cerr

	if err = cmd.Run(); err != nil {
		return err, errors.New(cerr.String())
	}

	return nil, nil
}

func DeleteBundleUser(username string) (error, error) {
	var cerr bytes.Buffer

	cmd := exec.Command("deluser", "--remove-all-files", username)
	cmd.Stderr = &cerr

	if err := cmd.Run(); err != nil {
		return err, errors.New(cerr.String())
	}

	return nil, nil
}
