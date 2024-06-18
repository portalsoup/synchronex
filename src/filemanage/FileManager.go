package filemanage

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func CopyFile(src, dest string, overwrite bool, owner string, group string) error {
	// Check if destination file already exists
	if isFilePresent(dest) {
		log.Printf("%s is already present\n", dest)
		if overwrite {
			err := DeleteFile(dest)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("destination file %s already exists\n", dest)
		}
	}

	// Open the source file
	srcFile, err := openFile(src)
	defer closeFile(srcFile)
	if err != nil {
		log.Printf("Failed to open src file %s\n", src)
		return err
	}

	// Create the destination file
	destFile, err := createFile(dest)
	defer closeFile(destFile)

	if err != nil {
		log.Printf("Failed to create dest file %s\n", dest)
		return err
	}
	// Copy the content from source to destination
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Set ownership of the new file
	foundUser, err := user.Lookup(owner)
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(foundUser.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(foundUser.Gid)
	if err != nil {
		return err
	}
	destFile.Chown(uid, gid)

	// Flush any buffered data to ensure file is written completely
	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func ValidateFileDoWork(path string, overwrite bool) bool {
	// Check if destination file already exists
	if isFilePresent(path) {
		return overwrite
	}
	return true

}

func FindChildren(rootPath string) ([]string, error) {
	var files []string
	suffix := ".nex.hcl"

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func openFile(file string) (*os.File, error) {
	log.Printf("Opening the file %s\n", file)
	openedFile, err := os.Open(file)
	if err != nil {
		return openedFile, err
	}
	return openedFile, nil
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func createFile(file string) (*os.File, error) {
	log.Printf("Creating the file at %s\n", file)

	// Ensure parent directories exist
	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		panic(err)
	}

	createdFile, err := os.Create(file)
	if err != nil {
		return createdFile, err
	}
	return createdFile, nil
}

func isFilePresent(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func DeleteFile(file string) error {
	log.Println("Deleting...")
	if !isFilePresent(file) {
		return fmt.Errorf("can't delete %s because it isn't present\n", file)
	}
	err := os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}
