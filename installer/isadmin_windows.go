//go:build windows
// +build windows

package main

import (
	"golang.org/x/sys/windows"
)

func isAdmin() bool {
	token, err := windows.OpenCurrentProcessToken()
	if err != nil {
		return false
	}
	defer token.Close()

	elevated := token.IsElevated()
	return elevated
}
