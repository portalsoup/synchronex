package src

import (
	"log"
	"path/filepath"
	"synchronex/src/filemanage"
	"synchronex/src/hcl"
	"synchronex/src/provision"
)

func ExecuteDocument(doc hcl.Document) {
	if doc.Sync {
		provision.Sync()
	}
	if doc.Upgrade {
		provision.Upgrade()
	}

	handlePackages(doc)
	handleScripts(doc)
	handleFiles(doc)

}

func handlePackages(doc hcl.Document) {
	for _, pkg := range doc.PackagesBlocks {
		switch pkg.Action {
		case "install":
			provision.Install(pkg.Package)
		case "remove":
			provision.Remove(pkg.Package)
		case "replace":
			{
				provision.Remove(pkg.Package)
				provision.Install(pkg.Package)
			}
		default:
		}
	}
}

func handleFiles(doc hcl.Document) {
	for _, file := range doc.FilesBlocks {
		sourceRaw, err := filepath.Abs(file.Source)
		if err != nil {
			log.Fatal(err)
		}

		source := hcl.Scan(doc, sourceRaw)
		dest := hcl.Scan(doc, file.Destination)

		switch file.Action {
		case "put":
			{
				err := filemanage.CopyFile(source, dest, false)
				if err != nil {
					log.Fatal(err)
				}
			}
		case "sync":
			{
				err := filemanage.CopyFile(source, dest, true)
				if err != nil {
					log.Fatal(err)
				}
			}
		case "remove":
			{
				log.Printf("About to delete %s\n", dest)
				err := filemanage.DeleteFile(dest)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func handleScripts(doc hcl.Document) {
	for _, script := range doc.ScriptsBlocks {
		switch script.Type {
		case "shell":
			runScript(script.ScriptSource)
		case "sh":
			runScript("/usr/bin/env", "sh", script.ScriptSource)
		case "bash":
			runScript("/usr/bin/env", "bash", script.ScriptSource)
		case "zsh":
			runScript("/usr/bin/env", "zsh", script.ScriptSource)
		case "fish":
			runScript("/usr/bin/env", "fish", script.ScriptSource)
		case "expect":
			runScript("/usr/bin/env", "expect", script.ScriptSource)
		default:
		}
	}
}

func runScript(src string, args ...string) {
	_, err := provision.Exec(src, args...)
	if err != nil {
		log.Fatal("failed to execute expect script")
	}
}
