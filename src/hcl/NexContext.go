package hcl

type NexContext struct {
	RequireRoot    bool   `hcl:"require_root,optional"`
	PersonalUser   string `hcl:"user,optional"`
	PackageManager string `hcl:"package_manager,optional"`
	Path           string
}
