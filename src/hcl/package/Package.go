package _package

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Package struct {
	// "put" copy if not present
	// "sync" unconditionally replace with current version
	// "remove" delete the file at the dest
	PackageManager string `hcl:"type,label"`
	Name           string `hcl:"name,label"`

	// If this file is to be copied, then it must have a source
	VersionCommand string `hcl:"version_command,optional"`
	VersionPattern string `hcl:"version_pattern,optional"`
	VersionRange   string `hcl:"range,optional"`
}

func (p Package) Executor() PackageExecutor {
	return PackageExecutor{
		Package: p,
	}
}

type PackageExecutor struct {
	Package Package
}

func (p PackageExecutor) Run() {
	p.Package.checkVersion()
}

func (p Package) checkVersion() {
	switch p.PackageManager {
	case "pacman":
		checkPacmanVersion(p)
	}
}

func checkPacmanVersion(p Package) {
	innerCommand := fmt.Sprintf("pacman -Qi %s | grep 'Version' | awk '{print $3}'", p.Name)
	cmd := exec.Command("bash", "-c", innerCommand)
	rawVersion, _ := cmd.Output()
	ranges, _ := TokenizeRange(p.VersionRange)

	version := strings.TrimSpace(string(rawVersion))
	for _, v := range ranges {
		log.Printf("Validating %s\n", p.Name)
		v.IsInRange(version)
	}
	log.Printf("%s was found and valid", p.Name)
}
