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
	for _, aRange := range ranges {
		log.Printf("Validating %s=%s", p.Name, version)
		if !aRange.IsInRange(version) {
			log.Println("[Error] Version mismatch!  Expected:")
			rangeVersionMismatchMessage(aRange, version)
		}
	}
}

func rangeVersionMismatchMessage(aRange Range, version string) {
	startOperand := "<"
	endOperand := ">"

	if aRange.Start.Inclusive {
		startOperand += "="
	}
	if aRange.End.Inclusive {
		endOperand += "="
	}

	logMsg := ""

	if aRange.Start.Version != "" {
		logMsg += aRange.Start.Version + " " + startOperand + " "
	}

	logMsg += version

	if aRange.End.Version != "" {
		logMsg += " " + endOperand + " " + aRange.End.Version
	}

	log.Fatalf("\t" + logMsg)
}
