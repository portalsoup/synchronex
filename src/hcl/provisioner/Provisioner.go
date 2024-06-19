package provisioner

import (
	"synchronex/src/hcl/file"
	"synchronex/src/hcl/folder"
	_package "synchronex/src/hcl/package"
)

type Provisioner struct {
	Name string `hcl:"type,label"`

	FilesBlocks    []file.File        `hcl:"file,block"`
	FoldersBlocks  []folder.Folder    `hcl:"folder,block"`
	PackagesBlocks []_package.Package `hcl:"package,block"`
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

func (p ProvisionExecutor) Run() {
	for _, packagesBlock := range p.Provisioner.PackagesBlocks {
		packagesBlock.Executor().Run()
	}

	for _, filesBlock := range p.Provisioner.FilesBlocks {
		filesBlock.Executor(p.User).Run()
	}

}
