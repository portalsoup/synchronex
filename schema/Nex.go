package schema

import (
	"log"
	"sort"
	"synchronex/common/execution"
	"synchronex/core"
)

type Nex struct {
	execution.Job[Nex] `json:"job,omitempty"`
	User               string `hcl:"user,optional" json:"user,omitempty"`
	Files              []File `hcl:"file,block" json:"files,omitempty"`
}

func (n *Nex) Validate() bool {
	return false
}

func (n *Nex) Execute() {

}

func (n *Nex) DifferencesFromState(state Nex) (diff *Nex) {
	diff = &Nex{
		Files:   make([]File, 0),
	}

	// First find files diff
	diff.Files = n.compareFiles(&state)

	if n.User != state.User {
		diff.User = n.User
	}

	return diff
}

func (n *Nex) compareFiles(state *Nex) (diff []File) {
	diff = make([]File, 0)
	userChanged := n.User != state.User

	plannedFiles := n.Files
	stateFiles := state.Files

	sort.Sort(FileSorter(plannedFiles))
	sort.Sort(FileSorter(stateFiles))

	// First compare the state with the plan to find things that no longer exist
	for _, stateFile := range stateFiles {
		if !core.Contains(plannedFiles, stateFile) {
			log.Println("Removing file: ", stateFile.Destination)
			diff = append(diff, File{
				Action:      "Remove",
				Destination: stateFile.Destination,
			})
		}
	}

	// Then compare the plan with the state to find out which things are new/changing
	for _, plannedFile := range plannedFiles {
		if core.Contains(stateFiles, plannedFile) && userChanged {
			log.Println("Updating file: ", plannedFile.Destination)
			diff = append(diff,
				File{
					Action:      "Remove",
					Destination: plannedFile.Destination,
				},
				File{
					Action:      "Add",
					Destination: plannedFile.Destination,
				},
			)
		} else if !core.Contains(stateFiles, plannedFile) {
			log.Println("Adding file: ", plannedFile.Destination)
			diff = append(diff, File{
				Action:      "Add",
				Destination: plannedFile.Destination,
			})
		}
	}
	return diff
}
