package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoBatchesDiff(t *testing.T) {
	// Input Data
	state := &Nex{
		Batches: []Nex{
			Nex{
				Files: []File{
					File{
						Source:      "test/file",
						Destination: "~/.cache/synchronex/file",
					},
				},
			},
		},
	}
	plan := &Nex{
		Batches: []Nex{
			Nex{
				Files: []File{
					File{
						Source:      "test/file",
						Destination: "~/.cache/synchronex/file",
					},
				},
			},
		},
	}

	// Expected Output
	expectedChanges := []Nex{} // Should have no changes

	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.Equal(t, expectedChanges, result.Batches)
}
