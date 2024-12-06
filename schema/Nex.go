package schema

import (
	"log"
	"sort"
	"synchronex/common/execution"
	"synchronex/common/hashcode"
)

type Nex struct {
	execution.Job[Nex] `json:"job,omitempty"`
	User               string `hcl:"user,optional" json:"user,omitempty"`
	Files              []File `hcl:"file,block" json:"files,omitempty"`
	Batches            []Nex  `hcl:"batch,block" json:"batches,omitempty"`
}

func (n *Nex) Validate() bool {
	return false
}

func (n *Nex) Execute() {

}

func (n *Nex) DifferencesFromState(state Nex) (diff *Nex) {
	diff = &Nex{
		User:    n.User,
		Files:   make([]File, 0),
		Batches: make([]Nex, 0),
	}
	userChanged := n.User != state.User

	plannedFiles := n.Files
	sort.Sort(FileSorter(plannedFiles))

	stateFiles := state.Files
	sort.Sort(FileSorter(stateFiles))

	// First compare the state with the plan to find things that no longer exist
	for _, stateFile := range stateFiles {
		if !contains(plannedFiles, stateFile) {
			log.Println("Removing file: ", stateFile.Destination)
			diff.Files = append(diff.Files, File{
				Action:      "Remove",
				Destination: stateFile.Destination,
				Source:      stateFile.Source,
				User:        stateFile.User,
				Group:       stateFile.Group,
				Permissions: stateFile.Permissions,
			})
		}
	}

	// Then compare the plan with the state to find out which things are new/changing
	for _, plannedFile := range plannedFiles {
		if contains(stateFiles, plannedFile) && userChanged {
			log.Println("Updating file: ", plannedFile.Destination)
			diff.Files = append(diff.Files, File{
				Action:      "Replace",
				Destination: plannedFile.Destination,
				Source:      plannedFile.Source,
				User:        plannedFile.User,
				Group:       plannedFile.Group,
				Permissions: plannedFile.Permissions,
			})
		} else if !contains(stateFiles, plannedFile) {
			log.Println("Adding file: ", plannedFile.Destination)
			diff.Files = append(diff.Files, File{
				Action:      "Add",
				Destination: plannedFile.Destination,
				Source:      plannedFile.Source,
				User:        plannedFile.User,
				Group:       plannedFile.Group,
				Permissions: plannedFile.Permissions,
			})
		}
	}

	return diff
}

func contains[T any](list []T, item T) bool {
	for _, v := range list {
		if hashcode.HashCode(v) == hashcode.HashCode(item) {
			return true
		}
	}
	return false
}
