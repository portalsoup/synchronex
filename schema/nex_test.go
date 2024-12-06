package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoDiff(t *testing.T) {
	// Input Data
	state := &Nex{
		User: "test",
		Files: []File{
			File{
				Source:      "test/file",
				Destination: "~/.cache/synchronex/file",
			},
		},
		Batches: []Nex{
			Nex{
				User: "test2",
				Files: []File{
					File{
						Source:      "test/file2",
						Destination: "~/.cache/synchronex/file2",
					},
				},
			},
		},
	}
	plan := &Nex{
		User: "test",
		Files: []File{
			File{
				Source:      "test/file",
				Destination: "~/.cache/synchronex/file",
			},
		},
		Batches: []Nex{
			Nex{
				User: "test2",
				Files: []File{
					File{
						Source:      "test/file2",
						Destination: "~/.cache/synchronex/file2",
					},
				},
			},
		},
	}

	// Expected Output
	expectedChanges := Nex{
		Files:   []File{},
		Batches: []Nex{},
	} // Should have no changes

	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.Equal(t, expectedChanges, *result)
}

func TestSetUser(t *testing.T) {
	// Input Data
	state := &Nex{}
	plan := &Nex{
		User: "test",
	}

	// Expected Output
	expectedChanges := Nex{
		User:    "test",
		Files:   []File{},
		Batches: []Nex{},
	} // Should have no changes

	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.Equal(t, expectedChanges, *result)
}
