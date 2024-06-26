package file

import (
	"log"
	"path/filepath"
	"synchronex/src/filemanage"
	"synchronex/src/hcl/context"
	"synchronex/src/provision"
	. "synchronex/src/shell"
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
	Shell       string `hcl:"shell,optional"`
	PreCommand  string `hcl:"pre_command,optional"`
	PostCommand string `hcl:"post_command,optional"`

	User  string `hcl:"owner,optional"`
	Group string `hcl:"group,optional"`
}

func (f File) Executor(context context.NexContext) FileExecutor {
	sourceRaw, err := filepath.Abs(f.Source)
	if err != nil {
		log.Fatal(err)
	}

	user := f.User
	if user == "" {
		user = context.PersonalUser
	}

	shellString := f.Shell
	if shellString == "" {
		shellString = "bash"
	}

	var shell Shell
	switch shellString {
	case "bash":
		{
			shell = Bash{}
		}
	case "zsh":
		{
			shell = Zsh{}
		}
	}

	templateConfig := template.Template{User: user}

	return FileExecutor{
		File:        f,
		Context:     context,
		Source:      templateConfig.Replace(sourceRaw),
		Destination: templateConfig.Replace(f.Destination),
		User:        user,
		Shell:       shell,
		Pre:         templateConfig.Replace(f.PreCommand),
		Post:        templateConfig.Replace(f.PostCommand),
	}
}

type FileExecutor struct {
	File        File
	Context     context.NexContext
	Source      string
	Destination string
	User        string
	Shell       Shell
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
	f.Shell.ExecuteCommand(f.Pre)

	// Do work
	success := copyFile(f, true)

	// Post script
	f.Shell.ExecuteCommand(f.Post)

	if !success {
		log.Printf("Failed       %s", f.Destination)
	} else {
		log.Printf("Synchronized %s", f.Destination)
	}
}

func (f FileExecutor) put() {
	// Pre script
	f.Shell.ExecuteCommand(f.Pre)

	// Do work
	success := copyFile(f, false)

	// Post script
	f.Shell.ExecuteCommand(f.Post)

	if !success {
		log.Printf("Failed       %s", f.Destination)
	} else {
		log.Printf("Copied       %s", f.Destination)
	}
}

func (f FileExecutor) remove() {
	// pre script
	f.Shell.ExecuteCommand(f.Pre)

	// Do work
	err := filemanage.DeleteFile(f.Destination)

	// post script
	f.Shell.ExecuteCommand(f.Post)

	if err != nil {
		log.Printf("Failed       %s", f.Destination)
	}
	log.Printf("Deleted      %s", f.Destination)

}
