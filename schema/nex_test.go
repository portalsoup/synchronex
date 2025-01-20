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
			{
				Source:      "test/file",
				Destination: "~/.cache/synchronex/file",
			},
		},
	}
	plan := &Nex{
		User: "test",
		Files: []File{
			{
				Source:      "test/file",
				Destination: "~/.cache/synchronex/file",
			},
		},
	}

	// Expected Output
	expectedChanges := &Nex{
		Files:   []File{},
	} // Should have no changes

	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.Equal(t, expectedChanges, result)
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
	} // Should have no changes

	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.Equal(t, expectedChanges, *result)
}
