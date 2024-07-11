package context

type NexContext struct {
	RequireRoot    bool   `hcl:"require_root,optional"`
	PersonalUser   string `hcl:"user"`
	PackageManager string `hcl:"package_manager,optional"`
}
