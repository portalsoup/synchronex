package common

import (
	"errors"
	"io"
	"log"
	"os"
	"path"
	"synchronex/schema"
)

// applyDiff applies the given list of file changes (diff).
func ApplyDiff(diff []schema.File) error {
	log.Printf("Applying %d changes...\n%s", len(diff), diff)
	for _, f := range diff {
		var err error
		switch f.Action {
		case "Add":
			err = applyAdd(f)
		case "Remove":
			err = applyRemove(f)
		default:
			log.Println("Unknown action:", f.Action)
			continue
		}

		// Log error if an operation fails
		if err != nil {
			log.Printf("Failed to apply action %s for %s: %v\n", f.Action, f.Destination, err)
			return err
		}
	}

	return nil
}

// applyAdd creates a file or folder at the given destination.
// If a source is provided, it copies the source content to the destination.
func applyAdd(f schema.File) error {
	if f.Source != "" {
		// If source is specified, copy the file to the destination
		srcFile, err := os.Open(f.Source)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Ensure target directory exists
		destDir := path.Dir(f.Destination)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return err
		}

		// Create the destination file
		destFile, err := os.Create(f.Destination)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Copy the source file content to the destination
		if _, err := io.Copy(destFile, srcFile); err != nil {
			return err
		}
		log.Println("File added:", f.Destination)
	} else {
		// If source is not specified, create an empty folder
		if err := os.MkdirAll(f.Destination, os.ModePerm); err != nil {
			return err
		}
		log.Println("Folder added:", f.Destination)
	}

	//if err := setFilePermissions(f); err != nil {
	//	return err
	//}

	return nil
}

// applyRemove deletes the specified file or folder.
func applyRemove(f schema.File) error {
	info, err := os.Stat(f.Destination)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("File or folder already removed:", f.Destination)
			return nil // The file/folder doesn't exist, nothing to do
		}
		return err
	}

	if info.IsDir() {
		// Remove the folder and its contents
		if err := os.RemoveAll(f.Destination); err != nil {
			return err
		}
		log.Println("Folder removed:", f.Destination)
	} else {
		// Remove the file
		if err := os.Remove(f.Destination); err != nil {
			return err
		}
		log.Println("File removed:", f.Destination)
	}

	return nil
}
