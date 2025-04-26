//go:build !windows
// +build !windows

package main

import "errors"

func installWindows() error {
	return errors.New("installWindows is not supported on this platform")
}
func uninstallWindows() error {
	return errors.New("uninstallWindows is not supported on this platform")
}
