package schema

import (
	"log"
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

	// If this file is to be copied, then it must have a source
	Source string `hcl:"src,optional" json:"source,omitempty"`
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
