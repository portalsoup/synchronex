package _package

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Package struct {
	Name string `hcl:"type,label"`

	// Supported package managers
	Pacman bool `hcl:"pacman,optional"`
	Dpkg   bool `hcl:"dpkg,optional"`

	VersionCommand string `hcl:"version_command,optional"`
	VersionPattern string `hcl:"version_pattern,optional"`
	VersionRange   string `hcl:"constraints,optional"`
}

func (p Package) Executor(pacmanInstalled, dpkgInstalled bool) PackageExecutor {
	return PackageExecutor{
		Package:        p,
		Pacman:         pacmanInstalled && p.Pacman,
		Dpkg:           dpkgInstalled && p.Dpkg,
		VersionCommand: p.VersionCommand,
		VersionPattern: p.VersionPattern,
		VersionRange:   p.VersionRange,
	}
}

type PackageExecutor struct {
	Package Package

	Pacman bool
	Dpkg   bool

	VersionCommand string
	VersionPattern string
	VersionRange   string
}

func (p PackageExecutor) Run() bool {
	return p.checkVersion()
}

func (p PackageExecutor) checkVersion() bool {
	var pacmanSuccess bool
	var aptSuccess bool
	if p.Pacman {
		pacmanSuccess = pacmanCheckVersion(p)
	}

	if p.Dpkg {
		aptSuccess = aptCheckVersion(p)
	}

	return pacmanSuccess || aptSuccess
}

func pacmanCheckVersion(p PackageExecutor) bool {
	cmdVersion := exec.Command("pacman", "-Qi", p.Package.Name)
	outputVersion, err := cmdVersion.Output()
	if cmdVersion.ProcessState.ExitCode() == 1 {
		log.Printf("Missing     %s", p.Package.Name)
		return false
	}

	// Extract the version information
	cmdGrep := exec.Command("grep", "Version")
	cmdGrep.Stdin = strings.NewReader(string(outputVersion))
	outputGrep, err := cmdGrep.Output()
	if err != nil {
		log.Printf("Error extracting version information: %v", err)
		return false
	}

	cmdAwk := exec.Command("awk", "{print $3}")
	cmdAwk.Stdin = strings.NewReader(string(outputGrep))
	rawVersion, err := cmdAwk.Output()
	if err != nil {
		log.Printf("Error extracting version number: %v", err)
		return false
	}

	ranges, _ := TokenizeRange(p.VersionRange)
	version := strings.TrimSpace(string(rawVersion))

	for _, aRange := range ranges {
		log.Printf("Validated   %s=%s", p.Package.Name, version)
		if !aRange.IsInRange(version) {
			log.Printf("Bad Version %s\t %s", p.Package.Name, rangeVersionMismatchMessage(aRange, version))
			return false
		}
	}
	return true
}

func aptCheckVersion(p PackageExecutor) bool {
	innerCommand := fmt.Sprintf("dpkg-query -W -f='${Version}' %s", p.Package.Name)
	cmd := exec.Command("bash", "-c", innerCommand)

	rawVersion, _ := cmd.Output()
	if cmd.ProcessState.ExitCode() == 1 {
		log.Printf("Missing     %s", p.Package.Name)
		return false
	}

	ranges, _ := TokenizeRange(p.VersionRange)

	version := strings.TrimSpace(string(rawVersion))
	for _, aRange := range ranges {
		log.Printf("Validated   %s=%s", p.Package.Name, version)
		if !aRange.IsInRange(version) {
			log.Printf("Bad Version %s --- %s", p.Package.Name, rangeVersionMismatchMessage(aRange, version))
			return false
		}
	}
	return true
}

func rangeVersionMismatchMessage(aRange Range, version string) string {
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

	return logMsg
}
