package schema

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"synchronex/common/execution"
	"synchronex/common/hashcode"
)

type FileSorter []File

func (a FileSorter) Len() int           { return len(a) }
func (a FileSorter) Less(i, j int) bool { return a[i].HashCode() < a[j].HashCode() }
func (a FileSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type File struct {
	execution.Job[File] `json:"job,omitempty"`

	Action      string `json:"action,omitempty"`
	Destination string `hcl:"type,label" json:"destination,omitempty"`
	Source      string `hcl:"src,optional" json:"source,omitempty"`
}

func (f *File) Validate() bool {
	return false
}

func (f *File) Execute() {

}

func (f *File) DifferencesFromState(state interface{}) interface{} {
	// Cast the state to the expected type ([]File)
	plannedFiles, ok := state.([]File)
	if !ok {
		// If the state is not of the expected type, return nil or an appropriate fallback
		log.Println("Invalid state type provided to DifferencesFromState")
		return nil
	}

	// Check if the current file exists in the planned files
	for _, plannedFile := range plannedFiles {
		if f.Destination == plannedFile.Destination {
			// File exists in the plan, no differences
			return nil
		}
	}

	// File does not exist in the planned files, return a difference
	return &File{
		Action:      "Remove",
		Destination: f.Destination,
	}
}

func (f *File) HashCode() uint32 {
	hash := hashcode.HashCode(f)
	return hash
}

func (f *File) ExpandFolderWildcard() ([]File, error) {
	suffix := "*"
	if strings.HasSuffix(f.Source, suffix) {
		err := validateDestinationDirectory(f.Destination)
		if err != nil {
			return nil, err
		}
		folderPath := strings.TrimSuffix(f.Source, suffix)

		var files []File

		entries, err := os.ReadDir(folderPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory: %w", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				srcPath := filepath.Join(folderPath, entry.Name())
				destPath := filepath.Join(f.Destination, entry.Name())

				files = append(files, File{
					Source:      srcPath,
					Destination: destPath,
				})
			}
		}

		return files, nil
	}

	return nil, nil
}

func validateDestinationDirectory(destination string) error {
	if stat, err := os.Stat(destination); err == nil {
		if !stat.IsDir() {
			return fmt.Errorf("destination '%s' exists and is not a directory", destination)
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(destination, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create destination directory '%s': %w", destination, err)
		}
	} else {
		return fmt.Errorf("failed to validate destination '%s': %w", destination, err)
	}
	return nil
}
