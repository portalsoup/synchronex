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

	summary := ""
	if nex, ok := v.(*schema.Nex); ok {
		add, remove := nex.DiffSummary()
		summary = fmt.Sprintf("\n\n%d to add, %d to remove", add, remove)
	}
	return fmt.Sprintf("%s%s", string(prettyJSON), summary)
}

func ParseNexFile(path string) (*schema.Nex, error) {
	var config schema.Nex
	err := hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		return nil, fmt.Errorf("Failed to load configuration: %s", err)
	}

	err = replaceWildcards(&config)
	if err != nil {
		return nil, fmt.Errorf("Failed to expand wildcards: %s", err)
	}
	return &config, nil
}

func WriteStatefile(state schema.Nex) (err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("error getting config dir: %v", err)
	}

	filePath := filepath.Join(configDir, "synchronex", "state.hcl")

	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating config dir: %v", err)
	}

	file, err := os.Create(filePath)
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
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("error getting config dir: %v", err)
	}

	filePath := filepath.Join(configDir, "synchronex", "state.hcl")

	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating config dir: %v", err)
	}

	config, err := ParseNexFile(filePath)
	if err != nil {
		return config, fmt.Errorf("Failed to load configuration: %s", err)
	}

	return config, nil
}

func replaceWildcards(n *schema.Nex) error {
	var expandedFiles []schema.File

	// Iterate over each File in the Nex struct
	for _, file := range n.Files {
		// Check and expand wildcards
		if expanded, err := file.ExpandFolderWildcard(); err != nil {
			return fmt.Errorf("failed to expand wildcard for file %s: %w", file.Source, err)
		} else if expanded != nil {
			// Append expanded files if wildcard was found
			expandedFiles = append(expandedFiles, expanded...)
		} else {
			// Append the original file if no wildcard was found
			expandedFiles = append(expandedFiles, file)
		}
	}

	// Replace the original Files array with the expanded list
	n.Files = expandedFiles
	return nil
}
