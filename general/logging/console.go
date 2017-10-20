// Package logging handles the various mediums and logic of logging messages and reports.
package logging

import (
	"log"
	"strings"

	"github.com/Ennovar/gPanel/general/networking"
)

const (
	PRIVATE_PREFIX string = "PRIVATE::"
	PUBLIC_PREFIX  string = "PUBLIC::"
)

const (
	NORMAL_LOG int = 1
	FATAL_LOG  int = 2
)

func Console(prefix string, logType int, msg string) {
	rawClientIP, _ := networking.GetClientIP()
	clientIP := strings.TrimSpace(string(rawClientIP))

	msg = prefix + clientIP + "::" + msg

	switch logType {
	default:
		fallthrough
	case 1:
		log.Println(msg)
	case 2:
		log.Fatal(msg)
	}
}
