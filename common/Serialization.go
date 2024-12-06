package common

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"log"
	"os"
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
	file, err := os.Create("statefile.hcl")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&state, f.Body())
	fmt.Printf("%s", f.Bytes())

	// Write the HCL bytes to the file
	if _, err := file.Write(f.Bytes()); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
