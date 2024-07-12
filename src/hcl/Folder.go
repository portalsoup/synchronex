package hcl

import (
	"log"
)

type Folder struct {
	// "put" copy if not present
	// "sync" unconditionally replace with current version
	// "remove" delete the file at the dest
	Action      string `hcl:"type,label"`
	Destination string `hcl:"name,label"`

	// If this file is to be copied, then it must have a source
	Source      string `hcl:"src,optional"`
	PreCommand  string `hcl:"pre_command,optional"`
	PostCommand string `hcl:"post_command,optional"`
}

func (f FolderExecutor) Validate() {
	log.Fatal("Folder validation not yet implemented!")
}

func (f Folder) Executor(context NexContext) FolderExecutor {
	return FolderExecutor{}
}

type FolderExecutor struct {
	Folder Folder
}
