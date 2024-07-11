package hcl

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ParseNexFile(path string) Nex {
	var config Nex
	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	return config
}

func ParseNexFiles(nexes []string) []Nex {
	foundNexes := make([]Nex, len(nexes))
	for i, nex := range nexes {
		config := ParseNexFile(nex)
		config.Context.Path = filepath.Dir(nex)
		foundNexes[i] = config
	}
	return foundNexes
}

func FindNexes(workingDir ...string) []string {
	var dir string
	if len(workingDir) > 0 {
		dir = workingDir[0]
	} else {
		// No working directory provided, use current working directory
		var err error
		dir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(os.Args) > 1 {
		return os.Args[1:]
	} else {
		foundNexes, err := getNexesInDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		return foundNexes
	}
}

func getNexesInDir(dir string) ([]string, error) {
	var nexes []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".nex.hcl") {
			nexes = append(nexes, filepath.Join(dir, entry.Name()))
		}
	}
	return nexes, nil
}
