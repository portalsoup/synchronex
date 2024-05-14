package schema

import (
	"fmt"
	"log"
	"os"
)

type Nex struct {
	Name         string `hcl:"name"`
	RequireRoot  bool   `hcl:"require_root,optional"`
	Sync         bool   `hcl:"sync_repositories,optional"`
	Upgrade      bool   `hcl:"upgrade_system,optional"`
	PersonalUser string `hcl:"user"`

	ProvisionerBlock Provisioner `hcl:"provisioner,block"`
}

func (n Nex) Executor() NexExecutor {
	name := n.Name
	if name == "" {
		name = n.ProvisionerBlock.Name
	}
	return NexExecutor{
		Nex:  n,
		Name: name,
		User: n.PersonalUser,
	}
}

func (n Nex) Validate() {

	err := validateRootRequirement(n)
	if err != nil {
		log.Fatal(err)
	}
}

// Validation functions

func validateRootRequirement(doc Nex) error {
	// If requires root, but is not root... fail
	if doc.RequireRoot == true && os.Geteuid() != 0 {
		return fmt.Errorf("The nex \"%s\" requires root\n", doc.Name)
	}
	return nil
}

type NexExecutor struct {
	Nex  Nex
	Name string
	User string
}

func (n NexExecutor) Run() {
	n.Nex.ProvisionerBlock.Executor(n.User).Run(n.Nex.Sync, n.Nex.Upgrade)
}
