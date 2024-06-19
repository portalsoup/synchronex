package file

import (
	"log"
	"path/filepath"
	"synchronex/src/filemanage"
	"synchronex/src/provision"
	"synchronex/src/template"
)

type File struct {
	// "put" copy if not present
	// "sync" unconditionally replace with current version
	// "remove" delete the file at the dest
	Action      string `hcl:"type,label"`
	Destination string `hcl:"name,label"`

	// If this file is to be copied, then it must have a source
	Source      string `hcl:"src,optional"`
	PreCommand  string `hcl:"pre_command,optional"`
	PostCommand string `hcl:"post_command,optional"`

	User  string `hcl:"owner,optional"`
	Group string `hcl:"group,optional"`
}

func (f File) Executor(defaultUser string) FileExecutor {
	sourceRaw, err := filepath.Abs(f.Source)
	if err != nil {
		log.Fatal(err)
	}

	user := f.User
	if user == "" {
		user = defaultUser
	}

	templateConfig := template.Template{User: user}

	return FileExecutor{
		File:        f,
		Source:      templateConfig.Replace(sourceRaw),
		Destination: templateConfig.Replace(f.Destination),
		User:        user,
		Pre:         templateConfig.Replace(f.PreCommand),
		Post:        templateConfig.Replace(f.PostCommand),
	}
}

type FileExecutor struct {
	File        File
	Source      string
	Destination string
	User        string
	Pre         string
	Post        string
}

func (f FileExecutor) Run() {
	switch f.File.Action {
	case "put":
		f.put()
	case "sync":
		f.sync()
	case "remove":
		f.remove()
	}
}

func (f FileExecutor) sync() {
	srcFile := f.Source
	destFile := f.Destination

	isEqual, err := provision.FilesEqual(srcFile, destFile)
	if err != nil {
		log.Fatal(err)
	}
	if isEqual {
		log.Printf("Up-To-Date   %s", f.Destination)
		return // The file is already present and up-to-date
	}

	// Pre script
	provision.ExecuteCommand(f.Pre)

	// Do work
	success := copyFile(f, true)

	// Post script
	provision.ExecuteCommand(f.Post)

	if !success {
		log.Printf("Failed       %s", f.Destination)
	} else {
		log.Printf("Synchronized %s", f.Destination)
	}
}

func (f FileExecutor) put() {
	// Pre script
	provision.ExecuteCommand(f.Pre)

	// Do work
	success := copyFile(f, false)

	// Post script
	provision.ExecuteCommand(f.Post)

	if !success {
		log.Printf("Failed       %s", f.Destination)
	} else {
		log.Printf("Copied       %s", f.Destination)
	}
}

func (f FileExecutor) remove() {
	// pre script
	provision.ExecuteCommand(f.Pre)

	// Do work
	err := filemanage.DeleteFile(f.Destination)

	// post script
	provision.ExecuteCommand(f.Post)

	if err != nil {
		log.Printf("Failed       %s", f.Destination)
	}
	log.Printf("Deleted      %s", f.Destination)

}
