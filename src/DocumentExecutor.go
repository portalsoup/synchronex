package src

import (
	"log"
	"path/filepath"
	"synchronex/src/filemanage"
	"synchronex/src/hcl"
	"synchronex/src/hcl/schema"
	"synchronex/src/provision"
)

func ExecuteDocument(doc schema.Provisioner) {
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

func handlePackages(doc schema.Provisioner) {
	for _, pkg := range doc.PackagesBlocks {
		switch pkg.Action {
		case "install":
			provision.Install(pkg, doc.PersonalUser, false)
		case "remove":
			provision.Remove(pkg)
		case "replace":
			{
				provision.Remove(pkg)
				provision.Install(pkg, doc.PersonalUser, false)
			}
		default:
		}
	}
}

func handleFiles(doc schema.Provisioner) {
	for _, file := range doc.FilesBlocks {
		sourceRaw, err := filepath.Abs(file.Source)
		if err != nil {
			log.Fatal(err)
		}

		source := hcl.Scan(doc, sourceRaw)
		dest := hcl.Scan(doc, file.Destination)

		switch file.Action {
		case "put":
			copyFile(file, source, dest, false)
		case "sync":
			copyFile(file, source, dest, true)
		case "remove":
			{
				log.Printf("About to delete %s\n", dest)
				// pre script?
				if filemanage.ValidateFileDoWork(dest, false) {
					_, err := provision.Exec("/usr/bin/bash", "-c", file.PreCommand)
					if err != nil {
						log.Fatal(err)
					}
				}
				err := filemanage.DeleteFile(dest)
				if err != nil {
					log.Fatal(err)
				}
				// post script?
				if filemanage.ValidateFileDoWork(dest, false) {
					_, err := provision.Exec("/usr/bin/bash", "-c", file.PostCommand)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}

func copyFile(file schema.File, source, dest string, overwrite bool) {
	// pre script?
	if filemanage.ValidateFileDoWork(dest, overwrite) && file.PreCommand != "" {
		log.Printf("Executing pre_command for %s", dest)
		_, err := provision.Exec("/usr/bin/bash", "-c", file.PreCommand)
		if err != nil {
			log.Fatal(err)
		}
	}

	// do work
	err := filemanage.CopyFile(source, dest, overwrite)
	if err != nil {
		log.Fatal(err)
	}

	// post script?
	if filemanage.ValidateFileDoWork(dest, overwrite) && file.PostCommand != "" {
		log.Printf("Executing post_command for %s", dest)
		_, err := provision.Exec("/usr/bin/bash", "-c", file.PostCommand)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handleScripts(doc schema.Provisioner) {
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

func runScript(args ...string) {
	_, err := provision.Exec(args...)
	if err != nil {
		log.Fatal("failed to execute expect script")
	}
}
