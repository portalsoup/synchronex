package schema

import (
	"log"
	"sort"
	"synchronex/common/execution"
	"synchronex/core"
)

type Nex struct {
	execution.Job[Nex] `json:"job,omitempty"`
	Files              []File `hcl:"file,block" json:"files,omitempty"`
}

func (n *Nex) Validate() bool {
	return false
}

func (n *Nex) Execute() {

}

func (n *Nex) DifferencesFromState(state Nex) (diff *Nex) {
	diff = &Nex{
		Files: make([]File, 0),
	}

	// First find files diff
	diff.Files = n.compareFiles(&state)

	return diff
}

func (n *Nex) compareFiles(state *Nex) (diff []File) {
	diff = make([]File, 0)

	plannedFiles := n.Files
	stateFiles := state.Files

	sort.Sort(FileSorter(plannedFiles))
	sort.Sort(FileSorter(stateFiles))

	// First compare the state with the plan to find things that no longer exist (Removals)
	for _, stateFile := range stateFiles {
		if difference := stateFile.DifferencesFromState(plannedFiles); difference != nil {
			// Cast the result back to a *File and add to the diff
			if diffFile, ok := difference.(*File); ok {
				diffFile.Action = "Remove"
				diff = append(diff, *diffFile)
			}
		}
	}

	// Then compare the plan with the state to find new files or changes (Additions and Updates)
	for _, plannedFile := range plannedFiles {
		log.Println("A planned file: ", plannedFile)
		if core.Contains(stateFiles, plannedFile) {
			// User has changed and the file exists in both states; mark it as "Update"
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
		} else if difference := plannedFile.DifferencesFromState(stateFiles); difference != nil {
			// Handle new files (Additions)
			if _, ok := difference.(*File); ok {
				diff = append(diff, File{
					Action:      "Add",
					Destination: plannedFile.Destination,
				})
			}
		}
	}
	return diff
}
