package src

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
	"os"
	"synchronex/src/hcl/schema"
)

func ProcessNexes(paths []string) {
	// Parse nexes paths into provisioners
	provisioners := readNexes(paths)

	errors := validateProvisioners(provisioners)
	if len(errors) > 0 {
		for _, err := range errors {
			log.Println(err)
		}
		log.Fatal("Exiting due to errors...")
	}

	// If they are validated successfully, then proceed to execute
	for _, provisioner := range provisioners {
		ExecuteDocument(provisioner)
	}
}

func IsRoot() bool {
	return os.Geteuid() == 0
}

func readNexes(nexes []string) []schema.Provisioner {
	configs := make([]schema.Provisioner, len(nexes))
	for i, nex := range nexes {
		config := readDocument(nex).ProvisionerBlock
		configs[i] = config
	}
	return configs
}

func validateProvisioners(provisioners []schema.Provisioner) []error {
	errors := make([]error, 0, len(provisioners))
	for _, provisioner := range provisioners {
		err := validateRootRequirement(provisioner)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func validateRootRequirement(doc schema.Provisioner) error {
	if doc.RequireRoot == true && IsRoot() != true {
		return fmt.Errorf("The nex \"%s\" requires root\n", doc.Name)
	}
	return nil
}

func readDocument(file string) schema.Document {
	var config schema.Document
	err := hclsimple.DecodeFile(file, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	return config
}
