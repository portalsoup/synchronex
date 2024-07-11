package hcl

import (
	"log"
	"os"
	"os/exec"
	"os/user"
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

func (f File) Validate() {
	// Validate shell exists
	if f.Shell != "" {
		if _, err := exec.LookPath(f.Shell); err != nil {
			log.Fatalf("shell %s is not installed", f.Shell)
		}
	}

	// Validate user exists
	if f.User != "" {
		if _, err := user.Lookup(f.User); err != nil {
			log.Fatalf("user %s does not exist", f.User)
		}
	}

	// Validate group exists
	if f.Group != "" {
		if _, err := user.LookupGroup(f.Group); err != nil {
			log.Fatalf("group %s does not exist", f.Group)
		}
	}

	// Validate source file exists at resource dir
	if f.Source != "" {
		if _, err := os.Stat(f.Source); os.IsNotExist(err) {
			log.Fatalf("source file %s does not exist", f.Source)
		}
	}

	// Validate destination dir is writable
	destDir := filepath.Dir(f.Destination)
	if err := isWritable(destDir); err != nil {
		log.Fatalf("destination directory %s is not writable: %v", destDir, err)
	}
}

func isWritable(path string) error {
	testFile := filepath.Join(path, ".writetest")
	file, err := os.OpenFile(testFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	file.Close()
	return os.Remove(testFile)
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
