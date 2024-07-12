package hcl

import (
	"fmt"
	"log"
	"os"
)

type Nex struct {
	Context *NexContext `hcl:"context,block"`

	ProvisionerBlock Provisioner `hcl:"provisioner,block"`
}

func (n Nex) Executor(context NexContext) NexExecutor {
	var newContext NexContext
	if n.Context != nil {
		newContext = *n.Context
	} else {
		newContext = context
	}

	return NexExecutor{
		Nex:     n,
		Name:    n.ProvisionerBlock.Name,
		Context: newContext,
	}
}

func (n NexExecutor) Validate() {

	err := validateRootRequirement(n.Nex)
	if err != nil {
		log.Fatal(err)
	}

	// propagate validation check
	n.Nex.ProvisionerBlock.Validate(n.Context)
}

// Validation functions

func validateRootRequirement(doc Nex) error {
	// If requires root, but is not root... fail
	if doc.Context != nil && doc.Context.RequireRoot == true && os.Geteuid() != 0 {
		return fmt.Errorf("The nex \"%s\" requires root\n", doc.ProvisionerBlock.Name)
	}
	return nil
}

type NexExecutor struct {
	Nex     Nex
	Context NexContext
	Name    string
	User    string
}

func (n NexExecutor) Run() {
	n.Nex.ProvisionerBlock.Executor(n.Context).Run()
}
