package main

import (
	"synchronex/src/hcl"
)

func main() {
	// Determine nexes finding strategy
	foundNexes := hcl.FindNexes()

	// Parse nexes rawPaths into objects
	nexes := hcl.ParseNexFiles(nil, foundNexes)

	// Validate each nex or fail
	for _, n := range nexes {
		n.Executor(hcl.NexContext{}).Validate()
	}

	// If they are validated successfully, then proceed to execute
	for _, n := range nexes {
		n.Executor(hcl.NexContext{}).Run()
	}
}
