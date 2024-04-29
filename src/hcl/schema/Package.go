package schema

type Package struct {
	Action  string `hcl:"type,label"`
	Package string `hcl:"name,label"`

	// If any package managers are specified, then this package will
	// be skipped for all other package managers.
	PackageManagers []string `hcl:"package_manager,optional"`
}
