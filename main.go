package main

import (
	"log"
	"os"
	"path/filepath"
	"synchronex/src/filemanage"
	"synchronex/src/hcl/schema"
)

func main() {
	// Determine nexes finding strategy
	foundNexes := findNexes()

	// Parse nexes rawPaths into objects
	nexes := schema.ParseNexFiles(foundNexes)

	// Validate each nex
	for _, nex := range nexes {
		nex.Validate()
	}
	// If they are validated successfully, then proceed to execute
	for _, nex := range nexes {
		nex.ProvisionerBlock.Handler().Run()
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
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	return filemanage.FindChildren(exPath)
}
