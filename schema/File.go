package schema

import (
	"synchronex/common/execution"
	"synchronex/common/hashcode"
)

type FileSorter []File

func (a FileSorter) Len() int           { return len(a) }
func (a FileSorter) Less(i, j int) bool { return a[i].HashCode() < a[j].HashCode() }
func (a FileSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type File struct {
	execution.Job[File] `json:"job,omitempty"`
	Action              string `json:"action,omitempty"`

	Destination string `hcl:"type,label" json:"destination,omitempty"`

	// If this file is to be copied, then it must have a source
	Source string `hcl:"src,optional" json:"source,omitempty"`

	User        string `hcl:"owner,optional" json:"user,omitempty"`
	Group       string `hcl:"group,optional" json:"group,omitempty"`
	Permissions string `hcl:"chmod,optional" json:"permissions,omitempty"`
}

func (f *File) Validate() bool {
	return false
}

func (f *File) Execute() {

}

func (f *File) DifferencesFromState(state interface{}) interface{} {
	return nil
}

func (f *File) HashCode() uint32 {
	hash := hashcode.HashCode(f)
	return hash
}
