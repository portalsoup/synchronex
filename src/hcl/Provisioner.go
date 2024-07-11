package hcl

import (
	"log"
	"synchronex/src/hcl/context"
)

type Provisioner struct {
	Name string `hcl:"type,label"`

	ModulesBlocks  []NexModule `hcl:"module,block"`
	FilesBlocks    []File      `hcl:"file,block"`
	FoldersBlocks  []Folder    `hcl:"folder,block"`
	PackagesBlocks []Package   `hcl:"package,block"`
}

func (p Provisioner) Executor(context context.NexContext) ProvisionExecutor {
	return ProvisionExecutor{
		Provisioner: p,
		Context:     context,
		User:        context.PersonalUser,
	}
}

func (p Provisioner) Validate() {
	for _, aModule := range p.ModulesBlocks {
		aModule.Validate()
	}

	for _, aFolder := range p.FoldersBlocks {
		aFolder.Validate()
	}

	for _, aFile := range p.FilesBlocks {
		aFile.Validate()
	}

	for _, aPackage := range p.PackagesBlocks {
		aPackage.Validate()
	}
}

type ProvisionExecutor struct {
	Provisioner Provisioner
	Context     context.NexContext
	User        string
}

func (p ProvisionExecutor) Run() {
	log.Println("***********************")
	log.Println("* Validating Packages *")
	log.Println("***********************")
	p.runPackages()

	log.Println("****************")
	log.Println("* Moving Files *")
	log.Println("****************")
	p.runFiles()

}

func (p ProvisionExecutor) runPackages() {
	failedPackages := false
	for _, packagesBlock := range p.Provisioner.PackagesBlocks {
		result := packagesBlock.Executor(p.Context).Run()
		if !result {
			failedPackages = true
		}
	}

	if failedPackages {
		log.Fatal("Some packages did not meet requirements!  Exiting.")
	}
}

func (p ProvisionExecutor) runFiles() {
	for _, filesBlock := range p.Provisioner.FilesBlocks {
		filesBlock.Executor(p.Context).Run()
	}
}
