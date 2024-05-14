package schema

import (
	"fmt"
	"log"
	"os"
)

type Nex struct {
	ProvisionerBlock Provisioner `hcl:"provisioner,block"`
}

func (n Nex) Validate() {

	err := validateRootRequirement(n.ProvisionerBlock)

	if err != nil {
		log.Fatal(err)
	}
}

// Validation functions

func validateRootRequirement(doc Provisioner) error {
	// If requires root, but is not root... fail
	if doc.RequireRoot == true && os.Geteuid() != 0 {
		return fmt.Errorf("The nex \"%s\" requires root\n", doc.Name)
	}
	return nil
}
