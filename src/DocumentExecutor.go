package src

import (
	"synchronex/src/hcl/schema"
	"synchronex/src/provision"
)

func ExecuteDocument(doc schema.Provisioner) {
	// System-level stuff
	if doc.Sync {
		provision.Sync()
	}
	if doc.Upgrade {
		provision.Upgrade()
	}

	// Packages
	for _, pkg := range doc.PackagesBlocks {
		pkg.Handler(doc.PersonalUser).Run()
	}

	// Scripts
	for _, script := range doc.ScriptsBlocks {
		script.Handler().Run()
	}

	// Files
	for _, file := range doc.FilesBlocks {
		file.Handler(doc.PersonalUser).Run()
	}
}
