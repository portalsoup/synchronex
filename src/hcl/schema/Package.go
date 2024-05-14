package schema

import "synchronex/src/provision"

type Package struct {
	PackageManager string `hcl:"type,label"`
	Package        string `hcl:"name,label"`

	// If any package managers are specified, then this package will
	// be skipped for all other package managers.
	Action string `hcl:"action"`
	AsUser string `hcl:"as_user,optional"`
}

func (p Package) Handler(defaultUser string) PackageExecutor {
	return PackageExecutor{
		User:       defaultUser,
		Pkg:        p,
		FailOnSkip: false,
	}

}

type PackageExecutor struct {
	User       string
	Pkg        Package
	FailOnSkip bool
}

func (p PackageExecutor) Run() {
	switch p.Pkg.Action {
	case "install":
		p.Install()
	case "remove":
		p.Remove()
	case "replace":
		p.Replace()
	default:
	}
}

func (p PackageExecutor) Install() {
	if provision.IsInstalled(p.Pkg, p.User) {
		provision.Install(p.Pkg, p.User, p.FailOnSkip)
	}
}

func (p PackageExecutor) Remove() {
	provision.Remove(p.Pkg)
}

func (p PackageExecutor) Replace() {
	provision.Remove(p.Pkg)
	provision.Install(p.Pkg, p.User, p.FailOnSkip)
}