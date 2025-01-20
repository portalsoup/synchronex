package schema

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"synchronex/common/execution"
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
		if difference := plannedFile.DifferencesFromState(stateFiles); difference != nil {
			// Handle new files (Additions)
			if _, ok := difference.(*File); ok {
				diff = append(diff, File{
					Action:      "Add",
					Destination: plannedFile.Destination,
					Source:      plannedFile.Source,
				})
			}
		}
	}
	return diff
}

func (n *Nex) DiffSummary() (add int, remove int) {
	// example output
	add = 0
	remove = 0

	for _, file := range n.Files {
		if file.Action == "Add" {
			add++
		} else if file.Action == "Remove" {
			remove++
		}
	}

	return add, remove

}
func (n *Nex) ExpandHomeFolder() {
	for i := range n.Files {
		n.Files[i].Destination = filepath.Join(os.Getenv("HOME"), n.Files[i].Destination)

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		n.Files[i].Source = filepath.Join(wd, n.Files[i].Source)
	}
}
