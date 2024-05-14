package provision

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"synchronex/src/hcl/schema"
)

func Install(pkg schema.Package, personalUser string, failOnSKip bool) {
	packageName := pkg.Package
	packageManager := pkg.PackageManager

	var selectedUser string
	if pkg.AsUser != "" {
		selectedUser = pkg.AsUser
	} else {
		selectedUser = personalUser
	}

	log.Printf("%s/%s slated for installation... checking requirements\n", packageName, packageManager)

	// Ensure selected package manager exists on the system
	if !isPackageManagerInstalled(packageManager) {
		msg := "%s %s because the %s is not installed"
		if failOnSKip {
			log.Fatalf(msg, "Failing", packageName, packageManager)
		} else {
			log.Printf(msg, "Skipping", packageName, packageManager)
			return
		}
	}

	if !install(packageName, packageManager, selectedUser) {
		log.Printf("%s installation failed\n", packageName)
	}
}

func IsInstalled(pkg schema.Package, personalUser string) bool {
	var selectedUser string
	if pkg.AsUser != "" {
		selectedUser = pkg.AsUser
	} else {
		selectedUser = personalUser
	}

	return isInstalled(pkg.Package, pkg.PackageManager, selectedUser)
}

func Remove(pkg schema.Package) {
	log.Printf("%s slated for removal... checking requirements\n", pkg)
	distroName, distroErr := distroName()

	if distroErr != nil {
		log.Fatal(distroErr)
	}

	packageName := pkg.Package

	if distroName == "Arch Linux" {
		log.Println("Found arch linux... using pacman")
		log.Printf("Removing package %s\n", pkg)
		if remove(packageName, pkg.PackageManager) {
			log.Printf("%s removed Successfully")
		} else {
			log.Printf("%s encountered an error during removal!")
		}
	}
}

func Upgrade() {
	log.Printf("Upgrading system...")
	distroName, distroErr := distroName()

	if distroErr != nil {
		log.Fatal(distroErr)
	}

	if distroName == "Arch Linux" {
		log.Println("Found arch linux... using pacman")
		log.Printf("Performing system upgrade")
		if pacmanUpdate() {
			log.Printf("Upgrade successful!")
		} else {
			log.Printf("Upgrade failed!")
		}
	}
}

func Sync() {
	log.Printf("Syncing system...")
	distroName, distroErr := distroName()

	if distroErr != nil {
		log.Fatal(distroErr)
	}

	if distroName == "Arch Linux" {
		log.Println("Found arch linux... using pacman")
		log.Printf("Syncing dependency repository")
		if pacmanSync() {
			log.Printf("Sync successful!")
		} else {
			log.Printf("Sync failed!")
		}
	}
}

func isPackageManagerInstalled(pkgManager string) bool {
	_, err := Exec("which", pkgManager)
	return err == nil
}

func install(pkg string, packageManager string, user string) bool {
	log.Printf("Installing as %s... %s\n", user, pkg)
	expectPath := fmt.Sprintf("../expect/%s/install.sh", packageManager)
	path, err := filepath.Abs(expectPath)
	if err != nil {
		log.Println(err)
	}

	// Execute expect script
	_, err = Exec("su", "-c", fmt.Sprintf("%s %s", path, pkg), user)

	if err != nil {
		log.Println(err)
	}

	return err == nil
}

func isInstalled(pkg string, pkgManager string, user string) bool {
	log.Printf("Checking if %s is installed by %s...", pkg, pkgManager)
	scriptPath := fmt.Sprintf("../expect/%s/check-installed.sh", pkgManager)

	path, err := filepath.Abs(scriptPath)
	if err != nil {
		log.Println(err)
	}

	_, err = Exec("su", "-c", fmt.Sprintf("%s %s", path, pkg), user)
	if err != nil {
		log.Println(err)
	}
	return err == nil
}

func remove(pkg string, packageManager string) bool {
	log.Printf("Removing... %s\n", pkg)
	log.Printf("Installing... %s\n", pkg)
	expectPath := fmt.Sprintf("../expect/%s/remove.sh", packageManager)
	path, err := filepath.Abs(expectPath)
	if err != nil {
		log.Println(err)
	}
	_, err = Exec(path, pkg)

	if err != nil {
		log.Println(err)
	}

	return err == nil
}

func pacmanSync() bool {
	_, err := Exec("pacman", "-Sy")

	if err != nil {
		log.Println(err)
	}

	return err == nil
}

func pacmanUpdate() bool {
	path, err := filepath.Abs("../expect/pacman/update.sh")
	if err != nil {
		log.Println(err)
	}
	_, err = Exec(path)

	if err != nil {
		log.Println(err)
	}

	return err == nil
}

func distroName() (string, error) {
	// Open the os-release file
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a map to store the key-value pairs
	distroInfo := make(map[string]string)

	// Read each line of the file and parse key-value pairs
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			// Remove double quotes if present
			value := strings.Trim(parts[1], "\"")
			distroInfo[parts[0]] = value
		}
	}

	// Check if distribution information exists
	if distroName, ok := distroInfo["NAME"]; ok {
		return distroName, nil
	}

	// Return an error if distribution information not found
	return "", errors.New("distribution information not found")
}
