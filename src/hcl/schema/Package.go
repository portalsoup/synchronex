package schema

type Package struct {
	PackageManager string `hcl:"type,label"`
	Package        string `hcl:"name,label"`

	// If any package managers are specified, then this package will
	// be skipped for all other package managers.
	Action string `hcl:"action"`
	AsUser string `hcl:"as_user,optional"`
}
