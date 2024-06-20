package nex

import (
	"fmt"
	"log"
	"os"
	"synchronex/src/hcl/context"
	"synchronex/src/hcl/provisioner"
)

type Nex struct {
	Context context.NexContext `hcl:"context,block"`

	ProvisionerBlock provisioner.Provisioner `hcl:"provisioner,block"`
}

func (n Nex) Executor(context context.NexContext) NexExecutor {
	name := n.ProvisionerBlock.Name

	return NexExecutor{
		Nex:     n,
		Name:    name,
		Context: n.Context,
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
	if doc.Context.RequireRoot == true && os.Geteuid() != 0 {
		return fmt.Errorf("The nex \"%s\" requires root\n", doc.ProvisionerBlock.Name)
	}
	return nil
}

type NexExecutor struct {
	Nex     Nex
	Context context.NexContext
	Name    string
	User    string
}

func (n NexExecutor) Run() {
	n.Nex.ProvisionerBlock.Executor(n.Context).Run()
}
