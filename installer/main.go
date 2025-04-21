package main

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed data/*
var dataFiles embed.FS

//go:embed install/*
var installFiles embed.FS

var (
	// These variables will be set at build time using ldflags
	AppName    string
	BinaryName string
)

const (
	choiceExit          = "0"
	choiceDefault       = "1"
	choicePortable      = "2"
	defaultInstallPerms = 0755
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("\nInstallation failed: %v\n", err)
		fmt.Println("\nPress any key to exit...")
		fmt.Scanln()
		os.Exit(1)
	}
	fmt.Println("Press any key to exit...")
	fmt.Scanln()
}

func run() error {
	fmt.Println("Welcome to the", AppName, "installer!")
	fmt.Println("Please choose an option:")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\n0. Exit")
		fmt.Println("1. Default installation (recommended)")
		fmt.Println("2. Portable installation")
		fmt.Print("Enter your choice: ")
		if !scanner.Scan() {
			fmt.Println("\nInput terminated, exiting installer...")
			return nil
		}
		choice := strings.TrimSpace(scanner.Text())

		var err error
		switch choice {
		case choiceExit:
			fmt.Println("Exiting installer...")
			return nil
		case choiceDefault:
			if runtime.GOOS == "linux" && !isRoot() {
				fmt.Println("This operation requires superuser privileges.\nPlease run the installer with sudo for a system-wide installation.")
				continue
			}
			fmt.Println("\nStarting default installation...")
			err = defaultInstallation()
		case choicePortable:
			fmt.Println("\nStarting portable installation...")
			err = portableInstallation()
		default:
			fmt.Println("Invalid choice:", choice)
			fmt.Println("Please choose 0, 1, or 2.")
			continue
		}
		if err != nil {
			return err
		}
		fmt.Println("\nInstallation completed successfully!")
		return nil
	}
}

func defaultInstallation() error {
	fmt.Println("Preparing installation directories...")
	configDir, err := os.UserConfigDir()
	if err != nil {
		return errors.New("failed to get user config directory: " + err.Error())
	}

	dataDest := filepath.Join(configDir, AppName)
	fmt.Printf("Creating configuration directory at: %s\n", dataDest)
	if err := os.MkdirAll(dataDest, defaultInstallPerms); err != nil {
		return errors.New("failed to create config directory: " + err.Error())
	}

	fmt.Println("Copying application data files...")
	if err := copyEmbeddedDir(dataFiles, "data", dataDest); err != nil {
		return errors.New("failed to copy data files: " + err.Error())
	}

	fmt.Println("Installing application for your platform...")
	switch runtime.GOOS {
	case "darwin":
		return installMacOS()
	case "linux":
		return installLinux()
	case "windows":
		return installWindows()
	default:
		return errors.New("unsupported platform: " + runtime.GOOS)
	}
}

func portableInstallation() error {
	fmt.Println("Setting up portable installation...")
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.New("failed to get current directory: " + err.Error())
	}

	appDir := filepath.Join(currentDir, AppName)
	dataDir := filepath.Join(appDir, "data")

	fmt.Printf("Creating application directory at: %s\n", appDir)
	if err := os.MkdirAll(dataDir, defaultInstallPerms); err != nil {
		return errors.New("failed to create portable directories: " + err.Error())
	}

	fmt.Println("Copying application data files...")
	if err := copyEmbeddedDir(dataFiles, "data", dataDir); err != nil {
		return errors.New("failed to copy data files: " + err.Error())
	}

	fmt.Println("Copying executable files...")
	switch runtime.GOOS {
	case "darwin":
		binaryPath := filepath.Join("install", AppName+".app", "Contents", "Resources", BinaryName)
		binaryDest := filepath.Join(appDir, BinaryName)
		return copyEmbeddedFile(installFiles, binaryPath, binaryDest)
	case "linux", "windows":
		binaryPath := "install" + "/" + BinaryName
		binaryDest := filepath.Join(appDir, BinaryName)
		return copyEmbeddedFile(installFiles, binaryPath, binaryDest)
	default:
		return errors.New("unsupported platform: " + runtime.GOOS)
	}
}

func installMacOS() error {
	fmt.Println("Installing application bundle for macOS...")
	appDest := filepath.Join("/Applications", AppName+".app")
	if err := copyEmbeddedDir(installFiles, "install", "/Applications"); err != nil {
		return errors.New("failed to copy .app bundle: " + err.Error())
	}

	fmt.Println("Setting application permissions...")
	return os.Chmod(appDest, defaultInstallPerms)
}

func installLinux() error {
	fmt.Println("Installing application for Linux...")
	binDest := filepath.Join("/usr/local/bin", BinaryName)
	fmt.Printf("Copying executable to: %s\n", binDest)
	if err := copyEmbeddedFile(installFiles, filepath.Join("install", BinaryName), binDest); err != nil {
		return errors.New("failed to copy binary: " + err.Error())
	}

	fmt.Println("Setting executable permissions...")
	if err := os.Chmod(binDest, defaultInstallPerms); err != nil {
		return errors.New("failed to set binary permissions: " + err.Error())
	}

	fmt.Println("Creating application menu entry...")
	desktopSrc := filepath.Join("install", AppName+".desktop")
	desktopData, err := installFiles.ReadFile(desktopSrc)
	if err != nil {
		return errors.New("failed to read desktop file: " + err.Error())
	}

	desktopContent := strings.ReplaceAll(string(desktopData), "$BinaryPath", binDest)
	desktopDest := filepath.Join("/usr/share/applications", AppName+".desktop")
	fmt.Printf("Creating desktop entry at: %s\n", desktopDest)
	if err := os.WriteFile(desktopDest, []byte(desktopContent), defaultInstallPerms); err != nil {
		return errors.New("failed to write desktop file: " + err.Error())
	}

	return nil
}

func copyEmbeddedDir(fs embed.FS, srcDir, destDir string) error {
	entries, err := fs.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := srcDir + "/" + entry.Name()
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(destPath, defaultInstallPerms); err != nil {
				return err
			}
			if err := copyEmbeddedDir(fs, srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyEmbeddedFile(fs, srcPath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyEmbeddedFile(fs embed.FS, srcPath, destPath string) error {
	data, err := fs.ReadFile(srcPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(destPath, data, defaultInstallPerms)
	if err != nil {
		return err
	}
	return nil
}

func isRoot() bool {
	return os.Getuid() == 0
}
