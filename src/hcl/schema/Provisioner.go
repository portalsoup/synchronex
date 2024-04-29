package schema

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
