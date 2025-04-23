//go:build !windows
// +build !windows

package main

func isAdmin() bool {
	return false
}
