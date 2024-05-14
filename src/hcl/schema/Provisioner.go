package schema

import "synchronex/src/provision"

type Provisioner struct {
	Name string `hcl:"type,label"`

	PackagesBlocks []Package `hcl:"package,block"`
	FilesBlocks    []File    `hcl:"file,block"`
	FoldersBlocks  []Folder  `hcl:"folder,block"`
	ScriptsBlocks  []Script  `hcl:"script,block"`

	RequireRoot  bool   `hcl:"require_root,optional"`
	Sync         bool   `hcl:"sync_repositories,optional"`
	Upgrade      bool   `hcl:"upgrade_system,optional"`
	PersonalUser string `hcl:"user"`
}

func (p Provisioner) Handler() ProvisionExecutor {
	return ProvisionExecutor{
		Provisioner: p,
	}
}

type ProvisionExecutor struct {
	Provisioner Provisioner
}

func (p ProvisionExecutor) Run() {
	// System-level stuff
	if p.Provisioner.Sync {
		provision.Sync()
	}
	if p.Provisioner.Upgrade {
		provision.Upgrade()
	}

	// Packages
	for _, pkg := range p.Provisioner.PackagesBlocks {
		pkg.Handler(p.Provisioner.PersonalUser).Run()
	}

	// Scripts
	for _, script := range p.Provisioner.ScriptsBlocks {
		script.Handler().Run()
	}

	// Files
	for _, file := range p.Provisioner.FilesBlocks {
		file.Handler(p.Provisioner.PersonalUser).Run()
	}
}
