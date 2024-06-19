package provisioner

import (
	"log"
	"os/exec"
	"synchronex/src/hcl/file"
	"synchronex/src/hcl/folder"
	_package "synchronex/src/hcl/package"
)

type Provisioner struct {
	Name string `hcl:"type,label"`

	FilesBlocks    []file.File        `hcl:"file,block"`
	FoldersBlocks  []folder.Folder    `hcl:"folder,block"`
	PackagesBlocks []_package.Package `hcl:"package,block"`
}

func (p Provisioner) Executor(user string) ProvisionExecutor {
	return ProvisionExecutor{
		Provisioner: p,
		User:        user,
	}
}

type ProvisionExecutor struct {
	Provisioner Provisioner
	User        string
}

func (p ProvisionExecutor) Run() {
	usePacman := isPacmanInstalled()
	useApt := isAptInstalled()

	if usePacman {
		log.Printf("Found pacman...")
	}
	if useApt {
		log.Printf("Found apt-get...")
	}
	for _, packagesBlock := range p.Provisioner.PackagesBlocks {
		packagesBlock.Executor(usePacman, useApt).Run()
	}

	for _, filesBlock := range p.Provisioner.FilesBlocks {
		filesBlock.Executor(p.User).Run()
	}

}

func isPacmanInstalled() bool {
	// Check if "pacman" command is available in PATH
	cmd := exec.Command("pacman", "--version")
	err := cmd.Run()
	return err == nil
}

func isAptInstalled() bool {
	// Check if "apt-get" command is available in PATH
	cmd := exec.Command("apt-get", "--version")
	err := cmd.Run()
	return err == nil
}
