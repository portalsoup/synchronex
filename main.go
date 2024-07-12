package main

import (
	"flag"
	"log"
	"synchronex/src/hcl"
)

func main() {

	plan := flag.Bool("plan", false, "Run validation without committing work")

	flag.Parse()

	filepaths := flag.Args()
	log.Println(filepaths)

	// Determine nexes finding strategy
	foundNexes := hcl.FindNexes(filepaths...)

	// Parse nexes rawPaths into objects
	nexes := hcl.ParseNexFiles(nil, foundNexes)

	// Validate each nex or fail
	log.Println("Running validations...")
	for _, n := range nexes {
		n.Executor(hcl.NexContext{}).Validate()
	}

	// If they are validated successfully and not a dryrun, then proceed to execute
	if !*plan {
		for _, n := range nexes {
			n.Executor(hcl.NexContext{
				Plan: *plan,
			}).Run()
		}
	} else {
		log.Println("Dry-run detected, skipping execution...")
	}
}
