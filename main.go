package main

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
	"os"
	"synchronex/src"
	"synchronex/src/hcl"
)

func main() {
	hclFile := readArgs()
	config := readDocument(hclFile)
	src.ExecuteDocument(config)
}

func readArgs() string {
	if len(os.Args) < 2 {
		log.Fatal("A provisioner file must be provided")
	}

	hclFile := os.Args[1]
	return hclFile
}

func readDocument(file string) hcl.Document {
	var config hcl.Document
	err := hclsimple.DecodeFile(file, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	log.Printf("Configuration is %#v", config)
	return config
}
