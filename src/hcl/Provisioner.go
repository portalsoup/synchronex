package hcl

import (
	"log"
)

type Provisioner struct {
	Name string `hcl:"type,label"`

	ModulesBlocks  []NexModule `hcl:"module,block"`
	FilesBlocks    []File      `hcl:"file,block"`
	FoldersBlocks  []Folder    `hcl:"folder,block"`
	PackagesBlocks []Package   `hcl:"package,block"`
}

func (p Provisioner) Executor(context NexContext) ProvisionExecutor {
	return ProvisionExecutor{
		Provisioner: p,
		Context:     context,
		User:        context.PersonalUser,
	}
}

func (p Provisioner) Validate(context NexContext) {
	for _, aModule := range p.ModulesBlocks {
		aModule.Executor(context).Validate()
	}

	for _, aFolder := range p.FoldersBlocks {
		aFolder.Executor(context).Validate()
	}

	for _, aFile := range p.FilesBlocks {
		aFile.Executor(context).Validate()
	}

	for _, aPackage := range p.PackagesBlocks {
		aPackage.Executor(context).Validate()
	}
}

type ProvisionExecutor struct {
	Provisioner Provisioner
	Path        string
	Context     NexContext
	User        string
}

func (p ProvisionExecutor) Run() {
	log.Println("***********************")
	log.Println("* Validating Packages *")
	log.Println("***********************")
	p.runPackages()

	log.Println("************************")
	log.Println("* Running nested Nexes *")
	log.Println("************************")
	p.runModules()

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

func (p ProvisionExecutor) runModules() {
	for _, moduleBlock := range p.Provisioner.ModulesBlocks {
		moduleBlock.Executor(p.Context).Run()
	}
}
