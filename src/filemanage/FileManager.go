package filemanage

import (
	"fmt"
	"io"
	"log"
	"os"
)

func CopyFile(src, dest string, overwrite bool) error {
	// Check if destination file already exists
	if isFilePresent(dest) {
		log.Println("Found the file already present")
		if overwrite {
			err := DeleteFile(dest)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("destination file %s already exists", dest)
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
		log.Printf("Failed to create dest file %s\n", src)
		return err
	}
	// Copy the content from source to destination
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Flush any buffered data to ensure file is written completely
	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
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
	log.Printf("Creating the file at %s", file)
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
