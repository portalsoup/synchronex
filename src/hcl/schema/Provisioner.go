package schema

import "synchronex/src/provision"

type Provisioner struct {
	Name string `hcl:"type,label"`

	PackagesBlocks []Package `hcl:"package,block"`
	FilesBlocks    []File    `hcl:"file,block"`
	FoldersBlocks  []Folder  `hcl:"folder,block"`
	ScriptsBlocks  []Script  `hcl:"script,block"`
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
	for _, script := range p.Provisioner.ScriptsBlocks {
		script.Executor().Run()
	}

	// Files
	for _, file := range p.Provisioner.FilesBlocks {
		file.Executor(p.User).Run()
	}
}
