package hcl

import (
	"log"
	"synchronex/src/filemanage"
)

func copyFile(file FileExecutor, overwrite bool) bool {
	// do work
	err := filemanage.CopyFile(file.Source, file.Destination, overwrite, file.User)
	if err != nil {
		log.Fatal(err)
	}

	return true
}
