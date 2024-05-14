package src

import (
	"fmt"
	"log"
	"os"
	"synchronex/src/hcl/schema"
)

func ProcessNexes(rawPaths []string) {
	// Parse nexes rawPaths into provisioners
	provisioners := readNexes(rawPaths)

	errors := validateProvisioners(provisioners)
	if len(errors) > 0 {
		for _, err := range errors {
			log.Println(err)
		}
		log.Fatal("Exiting due to errors...")
	}

	// If they are validated successfully, then proceed to execute
	for _, provisioner := range provisioners {
		provisioner.Handler().Run()
	}
}

func IsRoot() bool {
	return os.Geteuid() == 0
}

func readNexes(nexes []string) []schema.Provisioner {
	provisioners := make([]schema.Provisioner, len(nexes))
	for i, nex := range nexes {
		config := schema.ParseNexFile(nex).ProvisionerBlock
		provisioners[i] = config
	}
	return provisioners
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
