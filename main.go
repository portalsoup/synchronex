package main

import (
	"log"
	"synchronex/src/hcl"
)

func main() {
	// Determine nexes finding strategy
	foundNexes := hcl.FindNexes()

	// Parse nexes rawPaths into objects
	nexes := hcl.ParseNexFiles(nil, foundNexes)

	log.Printf("About to dereference a context")

	// Validate each nex or fail
	for _, n := range nexes {
		n.Executor(hcl.NexContext{}).Validate()
	}
	log.Printf("About to dereference a context")

	// If they are validated successfully, then proceed to execute
	for _, n := range nexes {
		n.Executor(hcl.NexContext{}).Run()
	}
}
