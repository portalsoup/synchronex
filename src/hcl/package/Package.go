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
	Apt    bool `hcl:"apt,optional"`

	VersionCommand string `hcl:"version_command,optional"`
	VersionPattern string `hcl:"version_pattern,optional"`
	VersionRange   string `hcl:"constraints,optional"`
}

func (p Package) Executor(usePacman, useApt bool) PackageExecutor {
	return PackageExecutor{
		Package:        p,
		Pacman:         usePacman,
		Apt:            useApt,
		VersionCommand: p.VersionCommand,
		VersionPattern: p.VersionPattern,
		VersionRange:   p.VersionRange,
	}
}

type PackageExecutor struct {
	Package Package

	Pacman bool
	Apt    bool

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

	if p.Apt {
		aptSuccess = aptCheckVersion(p)
	}

	return pacmanSuccess || aptSuccess
}

func checkPacmanVersion(p PackageExecutor) {
	// Check if the package is installed
	cmdCheck := exec.Command("pacman", "-Qi", p.Package.Name)
	outputCheck, err := cmdCheck.Output()
	if err != nil {
		log.Printf("Error running command or package not installed: %v", err)
		return
	}

	// Extract the version information
	cmdGrep := exec.Command("grep", "Version")
	cmdGrep.Stdin = strings.NewReader(string(outputCheck))
	outputGrep, err := cmdGrep.Output()
	if err != nil {
		log.Printf("Error extracting version information: %v", err)
		return
	}

	cmdAwk := exec.Command("awk", "{print $3}")
	cmdAwk.Stdin = strings.NewReader(string(outputGrep))
	rawVersion, err := cmdAwk.Output()
	if err != nil {
		log.Printf("Error extracting version number: %v", err)
		return
	}

	version := strings.TrimSpace(string(rawVersion))
	if version == "" {
		log.Printf("Package %s is not installed.", p.Package.Name)
		return
	}

	ranges, err := TokenizeRange(p.VersionRange)
	if err != nil {
		log.Printf("Error tokenizing version range: %v", err)
		return
	}

	log.Printf("Validating %s=%s", p.Package.Name, version)
	for _, aRange := range ranges {
		if !aRange.IsInRange(version) {
			log.Println("[Error] Version mismatch! Expected:")
			rangeVersionMismatchMessage(aRange, version)
		}
	}
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
