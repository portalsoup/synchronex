package hcl

import (
	"fmt"
	"log"
	"os"
	"synchronex/src/hcl/context"
)

type Nex struct {
	Context context.NexContext `hcl:"context,block"`

	ProvisionerBlock Provisioner `hcl:"provisioner,block"`

	Path string
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

	// propagate validation check
	n.ProvisionerBlock.Validate(n.Context)
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
