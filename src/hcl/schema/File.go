package schema

import (
	"log"
	"path/filepath"
	"synchronex/src/filemanage"
	"synchronex/src/hcl/template"
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

func (f File) Handler(defaultUser string) FileExecutor {
	sourceRaw, err := filepath.Abs(f.Source)
	if err != nil {
		log.Fatal(err)
	}

	user := f.User
	if user == "" {
		user = defaultUser
	}

	return FileExecutor{
		File:        f,
		Source:      template.ReplaceUser(user, sourceRaw),
		Destination: template.ReplaceUser(user, f.Destination),
		User:        user,
		Pre:         template.ReplaceUser(user, f.PreCommand),
		Post:        template.ReplaceUser(user, f.PostCommand),
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
		f.put(false)
	case "sync":
		f.put(true)
	case "remove":
		f.remove()
	}
}

func (f FileExecutor) put(overwrite bool) {
	log.Printf("About to write: %s", f.Destination)
	copied := copyFile(f.File, f.Source, f.Destination, overwrite, f.User)
	if copied {
		log.Printf(" ...done\n")
	} else {
		log.Printf(" ...skipped!\n")
	}
}

func (f FileExecutor) remove() {
	log.Printf("About to delete: %s", f.Destination)
	// pre script?
	executeBashCommand(f.Destination, false, f.Pre)

	// Do work
	err := filemanage.DeleteFile(f.Destination)
	if err != nil {
		log.Printf(" ...failed!\n")
		log.Fatal(err)
	}
	// post script?
	executeBashCommand(f.Destination, false, f.Post)
}
