package _package

type Package struct {
	// "put" copy if not present
	// "sync" unconditionally replace with current version
	// "remove" delete the file at the dest
	PackageManager string `hcl:"type,label"`
	Name           string `hcl:"name,label"`

	// If this file is to be copied, then it must have a source
	VersionCommand string `hcl:"version_command,optional"`
	VersionPattern string `hcl:"version_pattern,optional"`
	Versions       string `hcl:"versions,optional"`
}
