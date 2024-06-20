package shell

import (
	"log"
	"os/exec"
)

type Zsh struct{}

func (zsh Zsh) ExecuteCommand(cmd string) {
	_, err := Exec(false, whichZsh(), "-c", cmd)
	if err != nil {
		log.Fatal(err)
	}
}

func whichZsh() string {
	path, err := exec.LookPath("zsh")
	if err != nil {
		log.Fatal("Can't run command: zsh not found")
	}
	return path
}
