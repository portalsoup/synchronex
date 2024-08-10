package common

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"os"
	"path/filepath"
	"strings"
	"synchronex/schema"
)

func ParseNexFile(path string) (*schema.Nex, error) {
	var config schema.Nex
	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		return nil, fmt.Errorf("Failed to load configuration: %s", err)
	}
	return &config, nil
}

func GetNexes(searchDir string) ([]*schema.Nex, error) {
	var dir string
	if len(searchDir) > 0 {
		dir = searchDir
	} else {
		// No working directory provided, use current working directory
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("Failed to get current working directory: %s", err)
		}
	}

	foundNexes, err := findNexesInDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Failed to find Nexes: %s", err)
	}
	parsedNexes := make([]*schema.Nex, len(foundNexes))
	for i, nex := range foundNexes {
		aParsedNex, err := ParseNexFile(nex)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse Nex: %s", err)
		}
		parsedNexes[i] = aParsedNex
	}

	return parsedNexes, nil

}

func findNexesInDir(dir string) ([]string, error) {
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
