package main

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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
	choiceExit                 = "0"
	choiceSystemInstallation   = "1"
	choicePortableInstallation = "2"
	choiceSystemUninstallation = "3"
	systemInstallationOption   = "system_installation"
	systemUninstallationOption = "system_uninstallation"
)

var defaultInstallPerms os.FileMode = 0755

func main() {
	if err := run(); err != nil {
		fmt.Printf("\nOperation failed: %v\n", err)
		fmt.Println("\nPress any key to exit...")
		if _, err := fmt.Scanln(); err != nil {
			fmt.Printf("Input error: %v\n", err)
		}
		os.Exit(1)
	}
	fmt.Println("Press any key to exit...")
	if _, err := fmt.Scanln(); err != nil {
		fmt.Printf("Input error: %v\n", err)
	}
}

func run() error {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		var err error
		switch arg {
		case systemInstallationOption:
			if isUnix() && !isElevated() {
				fmt.Println("System installation requires administrator privileges.")
				fmt.Println("Please run with 'sudo' for system installation.")
				os.Exit(0)
			} else if runtime.GOOS == "windows" && !isElevated() {
				fmt.Println("System installation requires administrator privileges.")
				fmt.Println("Please run as Administrator for system installation.")
				os.Exit(0)
			}
			fmt.Println("Starting system installation...")
			err = systemInstallation()

		case systemUninstallationOption:
			if isUnix() && !isElevated() {
				fmt.Println("System removal requires administrator privileges.")
				fmt.Println("Please run with 'sudo' for system removal.")
				os.Exit(0)
			} else if runtime.GOOS == "windows" && !isElevated() {
				fmt.Println("System removal requires administrator privileges.")
				fmt.Println("Please run as Administrator for system removal.")
				os.Exit(0)
			}
			fmt.Println("Starting system removal...")
			err = systemUninstallation()
		}
		if err != nil {
			return err
		}
		fmt.Println("\nOperation completed successfully!")
		return nil
	}

	fmt.Println("Welcome to the", AppName, "installer!")
	fmt.Println("Please choose an option:")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\n0. Exit")
		fmt.Println("1. System installation (recommended)")
		fmt.Println("2. Portable installation")
		fmt.Println("3. Uninstall application")
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
		case choiceSystemInstallation:
			if !isElevated() {
				fmt.Println("System installation requires administrator privileges.")
				return restartWithElevatedPrivileges([]string{systemInstallationOption})
			}
			fmt.Println("\nStarting system installation...")
			err = systemInstallation()
		case choicePortableInstallation:
			fmt.Println("\nStarting portable installation in current directory...")
			err = portableInstallation()
		case choiceSystemUninstallation:
			if !isElevated() {
				fmt.Println("System removal requires administrator privileges.")
				return restartWithElevatedPrivileges([]string{systemUninstallationOption})
			}
			fmt.Println("\nStarting system removal...")
			err = systemUninstallation()
		default:
			fmt.Println("Invalid choice:", choice)
			fmt.Println("Please enter 0, 1, 2 or 3.")
			continue
		}

		if err != nil {
			return err
		}
		fmt.Println("\nOperation completed successfully!")
		return nil
	}
}

func restartWithElevatedPrivileges(args []string) error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	var cmdName string
	var cmdArgs []string

	switch runtime.GOOS {
	case "windows":
		cmdName = "powershell"
		quotedArgs := make([]string, len(args))
		for i, arg := range args {
			quotedArgs[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(arg, "'", "''"))
		}
		argsStr := strings.Join(quotedArgs, ", ")
		cmdArgs = []string{
			"-Command",
			fmt.Sprintf("Start-Process '%s' -ArgumentList %s -Verb runas", executable, argsStr),
		}
	case "darwin", "linux":
		cmdName = "sudo"
		cmdArgs = []string{
			executable,
		}
		cmdArgs = append(cmdArgs, args...)
	default:
		return errors.New("unsupported platform: " + runtime.GOOS)
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start elevated process: %v", err)
	}
	os.Exit(0)
	return nil
}

func systemInstallation() error {
	fmt.Println("Preparing installation directories...")
	configDir, err := OsConfigDir(runtime.GOOS)
	if err != nil {
		return errors.New("failed to get user config directory: " + err.Error())
	}

	dataDest := filepath.Join(configDir, AppName)
	fmt.Printf("Creating configuration directory at: %s\n", dataDest)
	if err := createDir(dataDest, defaultInstallPerms); err != nil {
		return errors.New("failed to create config directory: " + err.Error())
	}

	fmt.Println("Copying application data files...")
	if err := copyEmbeddedDir(dataFiles, "data", dataDest, defaultInstallPerms); err != nil {
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

func installMacOS() error {
	fmt.Println("Installing application bundle for macOS...")
	appDest := filepath.Join("/Applications", AppName+".app")
	if err := copyEmbeddedDir(installFiles, "install", "/Applications", defaultInstallPerms); err != nil {
		return errors.New("failed to copy .app bundle: " + err.Error())
	}

	fmt.Println("Setting application permissions...")
	if err := os.Chmod(appDest, defaultInstallPerms); err != nil {
		return errors.New("failed to set permissions: " + err.Error())
	}
	return nil
}

func installLinux() error {
	fmt.Println("Installing application for Linux...")
	binDest := filepath.Join("/usr/local/bin", BinaryName)
	fmt.Printf("Copying executable to: %s\n", binDest)
	if err := copyEmbeddedFile(installFiles, filepath.Join("install", BinaryName), binDest, defaultInstallPerms); err != nil {
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

func systemUninstallation() error {
	fmt.Println("Removing application configuration and data...")
	configDir, err := OsConfigDir(runtime.GOOS)
	if err != nil {
		return errors.New("failed to get user config directory: " + err.Error())
	}
	dataDir := filepath.Join(configDir, AppName)
	fmt.Printf("Removing configuration directory: %s\n", dataDir)
	if err := os.RemoveAll(dataDir); err != nil {
		return errors.New("failed to remove config directory: " + err.Error())
	}
	fmt.Println("Uninstalling application for your platform...")
	switch runtime.GOOS {
	case "darwin":
		return uninstallMacOS()
	case "linux":
		return uninstallLinux()
	case "windows":
		return uninstallWindows()
	default:
		return errors.New("unsupported platform: " + runtime.GOOS)
	}
}

func uninstallMacOS() error {
	fmt.Println("Uninstalling macOS application bundle...")
	appPath := filepath.Join("/Applications", AppName+".app")
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		fmt.Println("Application bundle not found, skipping...")
		return nil
	}
	fmt.Printf("Removing application bundle: %s\n", appPath)
	if err := os.RemoveAll(appPath); err != nil {
		return errors.New("failed to remove .app bundle: " + err.Error())
	}
	return nil
}

func uninstallLinux() error {
	fmt.Println("Uninstalling Linux application...")
	binPath := filepath.Join("/usr/local/bin", BinaryName)
	fmt.Printf("Removing executable: %s\n", binPath)
	if err := os.Remove(binPath); err != nil && !os.IsNotExist(err) {
		return errors.New("failed to remove binary: " + err.Error())
	}
	// Remove desktop entry
	desktopPath := filepath.Join("/usr/share/applications", AppName+".desktop")
	fmt.Printf("Removing desktop entry: %s\n", desktopPath)
	if err := os.Remove(desktopPath); err != nil && !os.IsNotExist(err) {
		return errors.New("failed to remove desktop entry: " + err.Error())
	}
	return nil
}

func portableInstallation() error {
	fmt.Println("Setting up portable installation...")
	currentDir, err := os.Getwd()
	if err != nil {
		return errors.New("failed to get current directory: " + err.Error())
	}

	appDir := filepath.Join(currentDir, AppName)
	dataDir := filepath.Join(appDir, "data")

	perms := defaultInstallPerms
	fmt.Printf("Creating application directory at: %s\n", appDir)
	if err := createDir(appDir, perms); err != nil {
		return errors.New("failed to create application directory: " + err.Error())
	}
	fmt.Printf("Creating data directory at: %s\n", dataDir)
	if err := createDir(dataDir, perms); err != nil {
		return errors.New("failed to create data directory: " + err.Error())
	}

	setOwnership := isRoot() && os.Getenv("SUDO_UID") != "" && os.Getenv("SUDO_GID") != ""

	fmt.Println("Copying application data files...")
	if err := copyEmbeddedDir(dataFiles, "data", dataDir, perms); err != nil {
		return errors.New("failed to copy data files: " + err.Error())
	}

	fmt.Println("Copying executable files...")
	var binaryDest string
	switch runtime.GOOS {
	case "darwin":
		binaryPath := filepath.Join("install", AppName+".app", "Contents", "Resources", BinaryName)
		binaryDest = filepath.Join(appDir, BinaryName)
		if err := copyEmbeddedFile(installFiles, binaryPath, binaryDest, perms); err != nil {
			return err
		}
	case "linux", "windows":
		binaryPath := "install" + "/" + BinaryName
		binaryDest = filepath.Join(appDir, BinaryName)
		if err := copyEmbeddedFile(installFiles, binaryPath, binaryDest, perms); err != nil {
			return err
		}
	default:
		return errors.New("unsupported platform: " + runtime.GOOS)
	}

	if setOwnership {
		if err := setOriginalUserOwnership(appDir); err != nil {
			fmt.Printf("Warning: failed to set ownership for %s: %v\n", appDir, err)
		}
	}

	return nil
}

func createDir(path string, perms os.FileMode) error {
	parentDir := filepath.Dir(path)
	if parentDir != "." && parentDir != "/" {
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			if err := createDir(parentDir, perms); err != nil {
				return err
			}
		}
	}
	if err := os.Mkdir(path, perms); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func copyEmbeddedDir(fs embed.FS, srcDir, destDir string, perms os.FileMode) error {
	parentDir := filepath.Dir(destDir)
	setOwnership := shouldSetOriginalUserOwnership(parentDir)

	if err := createDir(destDir, perms); err != nil {
		return err
	}
	if setOwnership && isRoot() {
		if err := setOriginalUserOwnership(destDir); err != nil {
			fmt.Printf("Warning: failed to set ownership for %s: %v\n", destDir, err)
		}
	}

	entries, err := fs.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := srcDir + "/" + entry.Name()
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			if err := copyEmbeddedDir(fs, srcPath, destPath, perms); err != nil {
				return err
			}
		} else {
			if err := copyEmbeddedFile(fs, srcPath, destPath, perms); err != nil {
				return err
			}
			if setOwnership && isRoot() {
				if err := setOriginalUserOwnership(destPath); err != nil {
					fmt.Printf("Warning: failed to set ownership for %s: %v\n", destPath, err)
				}
			}
		}
	}
	return nil
}

func copyEmbeddedFile(fs embed.FS, srcPath, destPath string, perms os.FileMode) error {
	data, err := fs.ReadFile(srcPath)
	if err != nil {
		return err
	}

	parentDir := filepath.Dir(destPath)
	setOwnership := shouldSetOriginalUserOwnership(parentDir)

	err = os.WriteFile(destPath, data, perms)
	if err != nil {
		return err
	}

	if setOwnership && isRoot() {
		if err := setOriginalUserOwnership(destPath); err != nil {
			fmt.Printf("Warning: failed to set ownership for %s: %v\n", destPath, err)
		}
	}
	return nil
}

func isRoot() bool {
	return os.Getuid() == 0
}

func OsConfigDir(platform string) (string, error) {
	switch platform {
	case "linux":
		return "/usr/local/share", nil
	case "windows":
		return "C:\\ProgramData", nil
	case "darwin":
		return "/Library/Application Support", nil
	default:
		return "", errors.New("Platform " + platform + " is not supported.")
	}
}

func setOriginalUserOwnership(path string) error {
	uidStr := os.Getenv("SUDO_UID")
	gidStr := os.Getenv("SUDO_GID")
	if uidStr == "" || gidStr == "" {
		return nil // Not running via sudo, no need to change ownership
	}

	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		return fmt.Errorf("invalid SUDO_UID: %v", err)
	}

	gid, err := strconv.Atoi(gidStr)
	if err != nil {
		return fmt.Errorf("invalid SUDO_GID: %v", err)
	}

	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chown(filePath, uid, gid)
	})
}

func shouldSetOriginalUserOwnership(parentDir string) bool {
	if !isRoot() {
		return false
	}

	uidStr := os.Getenv("SUDO_UID")
	gidStr := os.Getenv("SUDO_GID")
	if uidStr == "" || gidStr == "" {
		return false
	}

	canWrite, err := canUserWrite(parentDir)
	if err != nil {
		fmt.Printf("Warning: cannot check write permissions for %s: %v\n", parentDir, err)
		return false
	}

	return canWrite
}

func canUserWrite(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	mode := info.Mode().Perm()

	return mode&0022 != 0, nil // Право записи для группы (0020) или остальных (0002)
}
func isElevated() bool {
	switch runtime.GOOS {
	case "windows":
		return isAdmin()
	case "linux", "darwin":
		return isRoot()
	default:
		return false
	}
}
func isUnix() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin"
}
