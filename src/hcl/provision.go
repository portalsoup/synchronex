package hcl

type Document struct {
	PackagesBlocks []Package `hcl:"package,block"`
	FilesBlocks    []File    `hcl:"file,block"`
	ScriptsBlocks  []Script  `hcl:"script,block"`

	Sync         bool   `hcl:"sync_repositories,optional"`
	Upgrade      bool   `hcl:"upgrade_system,optional"`
	PersonalUser string `hcl:"user"`
}

type File struct {
	// "put" copy if not present
	// "sync" unconditionally replace with current version
	// "remove" delete the file at the dest
	Action      string `hcl:"type,label"`
	Destination string `hcl:"name,label"`

	// If this file is to be copied, then it must have a source
	Source string `hcl:"src,optional"`
}

type Package struct {
	Action  string `hcl:"type,label"`
	Package string `hcl:"name,label"`

	// If any package managers are specified, then this package will
	// be skipped for all other package managers.
	PackageManagers []string `hcl:"package_manager,optional"`
}

type Script struct {
	Type string `hcl:"type,label"`

	ScriptSource string `hcl:"src"`
}
