package hcl

import (
	"log"
	"os"
	"path/filepath"
	"synchronex/src/template"
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
	// Validate source directory exists
	srcInfo, err := os.Stat(f.Source)
	if os.IsNotExist(err) {
		log.Fatalf("source directory %s does not exist", f.Source)
	} else if !srcInfo.IsDir() {
		log.Fatalf("source path %s is not a directory", f.Source)
	}

	// Validate destination directory is a directory if it exists, and is writable
	destInfo, err := os.Stat(f.Dest)
	if err != nil {
		if !os.IsNotExist(err) { // if error is not because the directory does not exist
			log.Fatalf("failed to access destination directory: %v", err)
		} // if directory does not exist, no need to further check if it's a directory. It's expected to be created later.
	} else if !destInfo.IsDir() {
		log.Fatalf("destination path %s exists but is not a directory", f.Dest)
	} else {
		// Try to write to the destination path
		err = isWritable(f.Dest)
		if err != nil {
			log.Fatalf("failed to write to destination path %s: %v", f.Dest, err)
		}
	}
}

func (f Folder) Executor(context NexContext) FolderExecutor {
	tem := template.Template{User: context.PersonalUser}

	return FolderExecutor{
		Folder:  f,
		Context: context,
		Source:  tem.Replace(filepath.Join(context.Path, f.Source)),
		Dest:    tem.Replace(f.Destination),
		Pre:     tem.Replace(f.PreCommand),
		Post:    tem.Replace(f.PostCommand),
	}
}

type FolderExecutor struct {
	Folder  Folder
	Context NexContext
	Source  string
	Dest    string
	Pre     string
	Post    string
}

// Run method for the FolderExecutor that handles Actions: "put", "sync", "remove"
func (f FolderExecutor) Run() {
	srcFiles, err := os.ReadDir(f.Source)
	if err != nil {
		log.Fatalf("Unable to read source directory: %s", err)
	}
	for _, file := range srcFiles {
		if file.IsDir() {
			continue // ignoring subdirectories
		}

		fileHCL := File{
			Action:      f.Folder.Action,
			Destination: filepath.Join(f.Dest, file.Name()),
			Source:      filepath.Join(f.Source, file.Name()),
			PreCommand:  f.Pre,
			PostCommand: f.Post,
		}

		fileHCL.Executor(f.Context).Run()
	}
}
