package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
	"os"
	"path/filepath"
	"synchronex/src"
	"synchronex/src/filemanage"
	"synchronex/src/hcl"
)

func main() {
	log.Printf("Root check... %t\n", isRoot())
	if len(os.Args) > 2 {
		config := readDocument(os.Args[1])
		err := validateRootRequirement(config.ProvisionerBlock)
		if err != nil {
			log.Fatal(err)
		}
		src.ExecuteDocument(config.ProvisionerBlock)
	} else {
		nexes, err := findNexes()
		if err != nil {
			return
		}

		provisioners := readNexes(nexes)

		errors := validateProvisioners(provisioners)
		if len(errors) > 0 {
			for _, err := range errors {
				log.Println(err)
			}
			log.Fatal("Exiting due to errors...")
		}

		// If they are validated successfully, then proceed to execute
		for _, provisioner := range provisioners {
			src.ExecuteDocument(provisioner)
		}
	}
}

func readNexes(nexes []string) []hcl.Provisioner {
	configs := make([]hcl.Provisioner, len(nexes))
	for i, nex := range nexes {
		config := readDocument(nex).ProvisionerBlock
		configs[i] = config
	}
	return configs
}

func validateProvisioners(provisioners []hcl.Provisioner) []error {
	errors := make([]error, 0, len(provisioners))
	for _, provisioner := range provisioners {
		err := validateRootRequirement(provisioner)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func validateRootRequirement(doc hcl.Provisioner) error {
	if doc.RequireRoot == true && isRoot() != true {
		return fmt.Errorf("The nex \"%s\" requires root\n", doc.Name)
	}
	return nil
}

func isRoot() bool {
	return os.Geteuid() == 0
}

func findNexes() ([]string, error) {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	return filemanage.FindChildren(exPath)
}

func readDocument(file string) hcl.Document {
	var config hcl.Document
	err := hclsimple.DecodeFile(file, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	return config
}
