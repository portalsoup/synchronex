package schema

import (
	"log"
	"synchronex/src/provision"
)

func runScript(args ...string) {
	_, err := provision.Exec(args...)
	if err != nil {
		log.Fatal("failed to execute expect script")
	}
}
