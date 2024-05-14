package schema

import (
	"log"
	"path/filepath"
	"synchronex/src/filemanage"
	"synchronex/src/hcl/template"
	"synchronex/src/provision"
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

func (f File) Handler(defaultUser string) FileHandler {
	sourceRaw, err := filepath.Abs(f.Source)
	if err != nil {
		log.Fatal(err)
	}

	return FileHandler{
		File:        f,
		Source:      template.ReplaceUser(defaultUser, sourceRaw),
		Destination: template.ReplaceUser(defaultUser, f.Destination),
		BackupUser:  defaultUser,
	}
}

type FileHandler struct {
	File        File
	Source      string
	Destination string
	BackupUser  string
}

func (f FileHandler) Run() {
	switch f.File.Action {
	case "put":
		f.put(false)
	case "sync":
		f.put(true)
	case "remove":
		f.remove()
	}
}

func (f FileHandler) put(overwrite bool) {
	log.Printf("About to write: %s")
	copied := copyFile(f.File, f.Source, f.Destination, overwrite, f.BackupUser)
	if copied {
		log.Printf(" ...done\n")
	} else {
		log.Printf(" ...skipped!\n")
	}
}

func (f FileHandler) remove() {
	log.Printf("About to delete: %s", f.Destination)
	// pre script?
	if filemanage.ValidateFileDoWork(f.Destination, false) {
		_, err := provision.Exec("/usr/bin/bash", "-c", f.File.PreCommand)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := filemanage.DeleteFile(f.Destination)
	if err != nil {
		log.Printf(" ...failed!\n")
		log.Fatal(err)
	}
	// post script?
	runBashScript(f.Destination, false, f.File)
}
