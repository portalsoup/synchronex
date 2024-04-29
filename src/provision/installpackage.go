package provision

import (
	"bufio"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Install(pkg string) {
	log.Printf("%s slated for installation... checking requirements\n", pkg)
	distroName, distroErr := distroName()

	if distroErr != nil {
		log.Fatal(distroErr)
	}

	if distroName == "Arch Linux" {
		log.Println("Found arch linux... using pacman")
		if !pacmanFind(pkg) {
			log.Printf("%s not found... installing\n", pkg)
			if pacmanInstall(pkg) {
				log.Printf("%s installed successfully!\n", pkg)
			} else {
				log.Printf("%s installation failed\n", pkg)
			}
		} else {
			log.Printf("%s already installed, skipping...\n", pkg)
		}
	}

}

func Remove(pkg string) {
	log.Printf("%s slated for removal... checking requirements\n", pkg)
	distroName, distroErr := distroName()

	if distroErr != nil {
		log.Fatal(distroErr)
	}

	if distroName == "Arch Linux" {
		log.Println("Found arch linux... using pacman")
		log.Printf("Removing package %s\n", pkg)
		if pacmanRemove(pkg) {
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

func pacmanInstall(pkg string) bool {
	log.Printf("Installing... %s\n", pkg)
	path, err := filepath.Abs("../expect/pacman/install.expect")
	if err != nil {
		log.Println(err)
	}
	_, err = Exec(path, pkg)

	if err != nil {
		log.Println(err)
	}

	return err == nil
}

func pacmanRemove(pkg string) bool {
	log.Printf("Removing... %s\n", pkg)
	path, err := filepath.Abs("../expect/pacman/remove.expect")
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
	path, err := filepath.Abs("../expect/pacman/update.expect")
	if err != nil {
		log.Println(err)
	}
	_, err = Exec(path)

	if err != nil {
		log.Println(err)
	}

	return err == nil
}

func pacmanFind(pkg string) bool {
	_, err :=
		exec.Command("pacman", "-Qi", pkg).Output()

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
