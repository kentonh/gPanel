// Package logging handles the various mediums and logic of logging messages and reports.
package logging

import (
	"log"
	"strings"

	"github.com/Ennovar/gPanel/pkg/networking"
)

const (
	PRIVATE_PREFIX = "PRIVATE::"
	PUBLIC_PREFIX  = "PUBLIC::"

	NORMAL_LOG = 1
	FATAL_LOG  = 2
)

// Console logs a prefix, IP, and message all appeneded together to the console.
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
