package provisioner

import (
	"synchronex/src/hcl/file"
	"synchronex/src/hcl/folder"
)

type Provisioner struct {
	Name string `hcl:"type,label"`

	FilesBlocks   []file.File     `hcl:"file,block"`
	FoldersBlocks []folder.Folder `hcl:"folder,block"`
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
	// Files
	for _, filesBlock := range p.Provisioner.FilesBlocks {
		filesBlock.Executor(p.User).Run()
	}
}
