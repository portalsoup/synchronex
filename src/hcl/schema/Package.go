package schema

import "synchronex/src/provision"

type Package struct {
	PackageManager string `hcl:"type,label"`
	Package        string `hcl:"name,label"`

	Action string `hcl:"action"`
	AsUser string `hcl:"as_user,optional"`
}

func (p Package) Handler(defaultUser string) PackageExecutor {
	return PackageExecutor{
		Package:    p,
		User:       defaultUser,
		Pkg:        p,
		FailOnSkip: false,
	}

}

type PackageExecutor struct {
	Package    Package
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
	if provision.IsInstalled(p.Pkg.Package, p.Pkg.PackageManager, p.User) {

		user := p.Package.AsUser
		if user == "" {
			user = p.User
		}
		provision.Install(p.Pkg.Package, p.Pkg.PackageManager, user, p.FailOnSkip)
	}
}

func (p PackageExecutor) Remove() {
	provision.Remove(p.Pkg.Package, p.Pkg.PackageManager)
}

func (p PackageExecutor) Replace() {
	user := p.Package.AsUser
	if user == "" {
		user = p.User
	}

	provision.Remove(p.Pkg.Package, p.Pkg.PackageManager)
	provision.Install(p.Pkg.Package, p.Pkg.PackageManager, user, p.FailOnSkip)
}
