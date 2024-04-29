package main

import (
	"log"
	"os"
	"path/filepath"
	"synchronex/src"
	"synchronex/src/filemanage"
)

func main() {
	log.Printf("Root check... %t\n", src.IsRoot())
	var foundNexes []string
	var err error

	// Determine nexes finding strategy
	if len(os.Args) > 1 {
		foundNexes = os.Args[1:]
	} else {
		foundNexes, err = findNexesInWorkingDir()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Begin processing
	src.ProcessNexes(foundNexes)
}

func findNexesInWorkingDir() ([]string, error) {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	return filemanage.FindChildren(exPath)
}
