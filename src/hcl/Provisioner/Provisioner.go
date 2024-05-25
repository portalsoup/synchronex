package Provisioner

import (
	"synchronex/src/hcl/file"
	"synchronex/src/hcl/folder"
	"synchronex/src/hcl/package"
	"synchronex/src/hcl/script"
	"synchronex/src/provision"
)

type Provisioner struct {
	Name string `hcl:"type,label"`

	PackagesBlocks []_package.Package `hcl:"package,block"`
	FilesBlocks    []file.File        `hcl:"file,block"`
	FoldersBlocks  []folder.Folder    `hcl:"folder,block"`
	ScriptsBlocks  []script.Script    `hcl:"script,block"`
}

func (p Provisioner) Executor(user string) ProvisionExecutor {
	return ProvisionExecutor{
		Provisioner: p,
		User:        user,
	}
}

type ProvisionExecutor struct {
	Provisioner Provisioner
	User        string
}

func (p ProvisionExecutor) Run(sync, upgrade bool) {
	// System-level stuff
	if sync {
		provision.Sync()
	}
	if upgrade {
		provision.Upgrade()
	}

	// Packages
	for _, pkg := range p.Provisioner.PackagesBlocks {
		pkg.Executor(p.User).Run()
	}

	// Scripts
	for _, scriptsBlock := range p.Provisioner.ScriptsBlocks {
		scriptsBlock.Executor().Run()
	}

	// Files
	for _, filesBlock := range p.Provisioner.FilesBlocks {
		filesBlock.Executor(p.User).Run()
	}
}
