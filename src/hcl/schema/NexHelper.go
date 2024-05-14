package schema

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
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
		foundNexes[i] = config
	}
	return foundNexes
}
