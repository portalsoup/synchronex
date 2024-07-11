package main

import (
	"synchronex/src/hcl"
	"synchronex/src/hcl/context"
)

func main() {
	// Determine nexes finding strategy
	foundNexes := hcl.FindNexes()

	// Parse nexes rawPaths into objects
	nexes := hcl.ParseNexFiles(foundNexes)

	// Validate each nex or fail
	for _, n := range nexes {
		n.Validate()
	}
	// If they are validated successfully, then proceed to execute
	for _, n := range nexes {
		n.Executor(context.NexContext{}).Run()
	}
}
