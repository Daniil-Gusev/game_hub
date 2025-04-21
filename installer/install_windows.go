//go:build windows
// +build windows

package main

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func installWindows() error {
	programFiles, err := windows.KnownFolderPath(windows.FOLDERID_ProgramFiles, 0)
	if err != nil {
		return errors.New("failed to get Program Files path: " + err.Error())
	}

	appDir := filepath.Join(programFiles, AppName)
	if err := os.MkdirAll(appDir, defaultInstallPerms); err != nil {
		return errors.New("failed to create app directory: " + err.Error())
	}

	binaryDest := filepath.Join(appDir, BinaryName)
	if err := copyEmbeddedFile(installFiles, ("install" + "/" + BinaryName), binaryDest); err != nil {
		return errors.New("failed to copy binary: " + err.Error())
	}

	// Create shortcut in Start Menu
	startMenu, err := windows.KnownFolderPath(windows.FOLDERID_StartMenu, 0)
	if err != nil {
		return errors.New("failed to get Start Menu path: " + err.Error())
	}

	shortcutPath := filepath.Join(startMenu, "Programs", AppName, AppName+".lnk")
	if err := os.MkdirAll(filepath.Dir(shortcutPath), defaultInstallPerms); err != nil {
		return errors.New("failed to create Start Menu directory: " + err.Error())
	}

	return createWindowsShortcut(binaryDest, shortcutPath)
}

func createWindowsShortcut(targetPath, shortcutPath string) error {
	// Note: This is a simplified version. For production, you might want to use
	// a proper COM-based approach to create Windows shortcuts.
	// This is a placeholder for the actual shortcut creation logic.
	return nil
}
