package common

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"log"
	"os"
	"path/filepath"
	"synchronex/schema"
)

func PrintPretty(v interface{}) string {
	prettyJSON, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatalln("Failed to generate pretty JSON:", err)
	}

	return string(prettyJSON)
}

func ParseNexFile(path string) (*schema.Nex, error) {
	var config schema.Nex
	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		return nil, fmt.Errorf("Failed to load configuration: %s", err)
	}
	return &config, nil
}

func WriteStatefile(state schema.Nex) (err error) {
	// Create or open a file in the current working directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// Create the full path to the .local directory
	dataDir := filepath.Join(homeDir, ".local", "share", "synchronex")
	statefile := filepath.Join(dataDir, "statefile.hcl")

	// Ensure the directory exists
	err = os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	file, err := os.Create(statefile)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&state, f.Body())
	fmt.Printf("%s", f.Bytes())

	// Write the HCL bytes to the file
	if _, err := file.Write(f.Bytes()); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func ReadStatefile() (*schema.Nex, error) {
	var config = schema.Nex{}
	//var config schema.Nex
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataDir := filepath.Join(homeDir, ".local", "share", "synchronex")
	statefile := filepath.Join(dataDir, "statefile.hcl")

	err = hclsimple.DecodeFile(statefile, nil, &config)
	log.Println("Found state:\t", config)
	if err != nil {
		return &config, fmt.Errorf("Failed to load configuration: %s", err)
	}
	return &config, nil
}
