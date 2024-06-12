package main

import (
	"log"
	"os"
	"path/filepath"
	"synchronex/src/filemanage"
	"synchronex/src/hcl/nex"
)

func main() {
	// Determine nexes finding strategy
	foundNexes := findNexes()

	// Parse nexes rawPaths into objects
	nexes := nex.ParseNexFiles(foundNexes)

	// Validate each nex or fail
	for _, n := range nexes {
		n.Validate()
	}
	// If they are validated successfully, then proceed to execute
	for _, n := range nexes {
		n.Executor().Run()
	}
}

func findNexes() []string {
	if len(os.Args) > 1 {
		return os.Args[1:]
	} else {
		foundNexes, err := getNexesInWorkingDir()
		if err != nil {
			log.Fatal(err)
		}
		return foundNexes
	}
}

func getNexesInWorkingDir() ([]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Dir(dir)
	return filemanage.FindChildren(path)
}
