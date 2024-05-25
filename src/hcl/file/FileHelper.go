package file

import (
	"log"
	"synchronex/src/filemanage"
	"synchronex/src/provision"
)

func executeBashCommand(dest string, overwrite bool, command string) {
	if filemanage.ValidateFileDoWork(dest, overwrite) {
		_, err := provision.Exec("/usr/bin/bash", "-c", command)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func copyFile(file File, source, dest string, overwrite bool, personalUser string) bool {
	// pre script?
	shouldSkip := preCopy(dest, overwrite, file)
	if shouldSkip {
		return false
	}

	user, group := getUseGroup(file, personalUser)

	// do work
	err := filemanage.CopyFile(source, dest, overwrite, user, group)
	if err != nil {
		log.Fatal(err)
	}

	// post script?
	postCopy(dest, overwrite, file)

	return true
}

func postCopy(dest string, overwrite bool, file File) {
	if filemanage.ValidateFileDoWork(dest, overwrite) && file.PostCommand != "" {
		log.Printf("Executing post_command for %s", dest)
		_, err := provision.Exec("/usr/bin/bash", "-c", file.PostCommand)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func preCopy(dest string, overwrite bool, file File) bool {
	if filemanage.ValidateFileDoWork(dest, overwrite) && file.PreCommand != "" {
		log.Printf("Executing pre_command for %s", dest)
		_, err := provision.Exec("/usr/bin/bash", "-c", file.PreCommand)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		return false // Skipping
	}
	return true
}

func getUseGroup(file File, backupUser string) (string, string) {
	user := ""
	if file.User != "" {
		user = file.User
	} else {
		user = backupUser
	}

	group := ""
	if file.Group != "" {
		group = file.Group
	} else {
		group = backupUser
	}
	return user, group
}
