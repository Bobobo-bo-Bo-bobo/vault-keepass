package main

import (
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"syscall"
)

func readStandardInput() string {
	raw, _ := terminal.ReadPassword(int(syscall.Stdin))
	return strings.Replace(strings.Replace(strings.Replace(string(raw), "\r", "", -1), "\n", "", -1), "\t", "", -1)
}
