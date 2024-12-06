package main

import (
	"github.com/stretchr/testify/assert"
	"synchronex/schema"
	"testing"
)

func TestNoFilesDiff(t *testing.T) {
	// Input Data
	state := &schema.Nex{
		Files: []schema.File{
			schema.File{
				Source:      "test/file",
				Destination: "~/.cache/synchronex/file",
			},
		},
	}
	plan := &schema.Nex{
		Files: []schema.File{
			schema.File{
				Source:      "test/file",
				Destination: "~/.cache/synchronex/file",
			},
		},
	}

	// Expected Output
	expectedChanges := []schema.File{} // Should have no changes

	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.Equal(t, expectedChanges, result.Files)
}

func TestAddFile(t *testing.T) {
	// Input Data
	file1 := schema.File{
		Source:      "test/file",
		Destination: "~/.cache/synchronex/file",
	}
	file2 := schema.File{
		Source:      "test/file2",
		Destination: "~/.cache/synchronex/file2",
	}

	state := &schema.Nex{
		Files: []schema.File{
			file1,
		},
	}
	plan := &schema.Nex{
		Files: []schema.File{
			file1,
			file2,
		},
	}

	// Expected Output
	expectedFile := schema.File{
		Action:      "Add",
		Source:      "test/file2",
		Destination: "~/.cache/synchronex/file2",
	}
	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.True(t, len(result.Files) == 1)
	assert.Equal(t, result.Files[0], expectedFile)
}

func TestRemoveFile(t *testing.T) {
	// Input Data
	file1 := schema.File{
		Source:      "test/file",
		Destination: "~/.cache/synchronex/file",
	}
	file2 := schema.File{
		Source:      "test/file2",
		Destination: "~/.cache/synchronex/file2",
	}

	state := &schema.Nex{
		Files: []schema.File{
			file1,
			file2,
		},
	}
	plan := &schema.Nex{
		Files: []schema.File{
			file1,
		},
	}

	// Expected Output
	expectedFile := schema.File{
		Action:      "Remove",
		Source:      "test/file2",
		Destination: "~/.cache/synchronex/file2",
	}
	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.True(t, len(result.Files) == 1)
	assert.Equal(t, result.Files[0], expectedFile)
}

func TestReplaceFile(t *testing.T) {
	// Input Data
	file1 := schema.File{
		Source:      "test/file",
		Destination: "~/.cache/synchronex/file",
	}
	newFile1 := schema.File{
		Source:      "test/file",
		Destination: "~/.cache/synchronex/file",
		User:        "test",
	}

	state := &schema.Nex{
		Files: []schema.File{
			file1,
		},
	}
	plan := &schema.Nex{
		Files: []schema.File{
			newFile1,
		},
	}

	// Expected Output
	expectedFiles := []schema.File{
		schema.File{
			Action:      "Remove",
			Source:      "test/file",
			Destination: "~/.cache/synchronex/file",
		},
		schema.File{
			Action:      "Add",
			Source:      "test/file",
			Destination: "~/.cache/synchronex/file",
			User:        "test",
		},
	}

	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.True(t, len(result.Files) == 2)
	assert.Equal(t, result.Files, expectedFiles)
}
