package emailer

import (
	"errors"
	"net/smtp"
	"strconv"
)

type Credentials struct {
	Username string
	Password string
	Server   string
	Port     int
}

type Emailer struct {
	auth smtp.Auth
	cred Credentials
}

func New(eType string, creds Credentials) (*Emailer, error) {
	var a smtp.Auth

	switch eType {
	case "crammd5":
		a = smtp.CRAMMD5Auth(creds.Username, creds.Password)
	default:
		a = smtp.PlainAuth("", creds.Username, creds.Password, creds.Server)
	}

	if a == nil {
		return nil, errors.New("unable to authenticate")
	}

	return &Emailer{
		auth: a,
		cred: creds,
	}, nil
}

// SendSimple function will fill out the to/subject email headers for you and allow
// you to just input a string as the email body.
func (e *Emailer) SendSimple(to string, subject string, body string) error {
	if e.auth == nil {
		return errors.New("smtp server authentication has expired")
	}

	m := []byte("To:" + to + "\r\n" +
		"Subject:" + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(e.cred.Server+":"+strconv.Itoa(e.cred.Port), e.auth, e.cred.Username, []string{to}, m)
	if err != nil {
		return err
	}

	return nil
}

// SendCustom function will not fill out the to/subject email headers for you and allow
// you to send a completely custom message in the form of bytes.
func (e *Emailer) SendCustom(to string, msg []byte) error {
	if e.auth == nil {
		return errors.New("smtp server authentication has expired")
	}

	err := smtp.SendMail(e.cred.Server+":"+strconv.Itoa(e.cred.Port), e.auth, e.cred.Username, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}
