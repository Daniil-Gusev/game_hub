//go:build windows
// +build windows

package main

import (
	"errors"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
)

func installWindows() error {
	fmt.Println("Installing application for Windows...")
	programFiles, err := windows.KnownFolderPath(windows.FOLDERID_ProgramFiles, 0)
	if err != nil {
		return errors.New("failed to get Program Files path: " + err.Error())
	}
	appDir := filepath.Join(programFiles, AppName)
	fmt.Printf("Creating application directory at: %s\n", appDir)
	if err := os.MkdirAll(appDir, defaultInstallPerms); err != nil {
		return errors.New("failed to create app directory: " + err.Error())
	}
	binaryDest := filepath.Join(appDir, BinaryName)
	fmt.Printf("Copying executable to: %s\n", binaryDest)
	if err := copyEmbeddedFile(installFiles, ("install" + "/" + BinaryName), binaryDest, defaultInstallPerms); err != nil {
		return errors.New("failed to copy binary: " + err.Error())
	}
	fmt.Println("Creating Start Menu shortcut...")
	startMenu, err := windows.KnownFolderPath(windows.FOLDERID_StartMenu, 0)
	if err != nil {
		return errors.New("failed to get Start Menu path: " + err.Error())
	}
	shortcutPath := filepath.Join(startMenu, "Programs", AppName, AppName+".lnk")
	fmt.Printf("Creating shortcut at: %s\n", shortcutPath)
	if err := os.MkdirAll(filepath.Dir(shortcutPath), defaultInstallPerms); err != nil {
		return errors.New("failed to create Start Menu directory: " + err.Error())
	}
	if err := createWindowsShortcut(binaryDest, shortcutPath); err != nil {
		return errors.New("failed to create shortcut: " + err.Error())
	}
	return nil
}
func createWindowsShortcut(targetPath, shortcutPath string) error {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return errors.New("failed to create WScript.Shell object: " + err.Error())
	}
	defer unknown.Release()

	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return errors.New("failed to query WScript.Shell interface: " + err.Error())
	}
	defer shell.Release()

	cs, err := oleutil.CallMethod(shell, "CreateShortcut", shortcutPath)
	if err != nil {
		return errors.New("failed to create shortcut: " + err.Error())
	}
	shortcut := cs.ToIDispatch()
	defer shortcut.Release()

	if _, err := oleutil.PutProperty(shortcut, "TargetPath", targetPath); err != nil {
		return errors.New("failed to set TargetPath: " + err.Error())
	}

	workingDir := filepath.Dir(targetPath)
	if _, err := oleutil.PutProperty(shortcut, "WorkingDirectory", workingDir); err != nil {
		return errors.New("failed to set WorkingDirectory: " + err.Error())
	}
	if _, err := oleutil.PutProperty(shortcut, "Description", "Shortcut for "+AppName); err != nil {
		return errors.New("failed to set Description: " + err.Error())
	}
	if _, err := oleutil.PutProperty(shortcut, "IconLocation", targetPath+",0"); err != nil {
		return errors.New("failed to set IconLocation: " + err.Error())
	}
	if _, err := oleutil.CallMethod(shortcut, "Save"); err != nil {
		return errors.New("failed to save shortcut: " + err.Error())
	}
	return nil
}

func uninstallWindows() error {
	fmt.Println("Removing Windows application...")
	programFiles, err := windows.KnownFolderPath(windows.FOLDERID_ProgramFiles, 0)
	if err != nil {
		return errors.New("failed to get Program Files path: " + err.Error())
	}
	appDir := filepath.Join(programFiles, AppName)
	fmt.Printf("Removing application directory: %s\n", appDir)
	if err := os.RemoveAll(appDir); err != nil && !os.IsNotExist(err) {
		return errors.New("failed to remove app directory: " + err.Error())
	}
	startMenu, err := windows.KnownFolderPath(windows.FOLDERID_StartMenu, 0)
	if err != nil {
		return errors.New("failed to get Start Menu path: " + err.Error())
	}
	shortcutDir := filepath.Join(startMenu, "Programs", AppName)
	shortcutPath := filepath.Join(shortcutDir, AppName+".lnk")
	fmt.Printf("Removing Start Menu shortcut: %s\n", shortcutPath)
	if err := os.Remove(shortcutPath); err != nil && !os.IsNotExist(err) {
		return errors.New("failed to remove shortcut: " + err.Error())
	}
	fmt.Printf("Removing Start Menu folder: %s\n", shortcutDir)
	if err := os.Remove(shortcutDir); err != nil && !os.IsNotExist(err) {
		return errors.New("failed to remove Start Menu folder: " + err.Error())
	}
	return nil
}
