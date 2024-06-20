package shell

import (
	"log"
	"os/exec"
)

type Bash struct {
}

func (bash Bash) ExecuteCommand(cmd string) {
	_, err := Exec(false, whichBash(), "-c", cmd)
	if err != nil {
		log.Fatal(err)
	}
}

func whichBash() string {
	path, err := exec.LookPath("bash")
	if err != nil {
		log.Fatal("Can't run command: bash not found")
	}
	return path
}
