package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoFilesDiff(t *testing.T) {
	// Input Data
	state := &Nex{
		Files: []File{
			File{
				Source:      "test/file",
				Destination: "~/.cache/synchronex/file",
			},
		},
	}
	plan := &Nex{
		Files: []File{
			File{
				Source:      "test/file",
				Destination: "~/.cache/synchronex/file",
			},
		},
	}

	// Expected Output
	expectedChanges := []File{} // Should have no changes

	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.Equal(t, expectedChanges, result.Files)
}

func TestAddFile(t *testing.T) {
	// Input Data
	file1 := File{
		Source:      "test/file",
		Destination: "~/.cache/synchronex/file",
	}
	file2 := File{
		Source:      "test/file2",
		Destination: "~/.cache/synchronex/file2",
	}

	state := &Nex{
		Files: []File{
			file1,
		},
	}
	plan := &Nex{
		Files: []File{
			file1,
			file2,
		},
	}

	// Expected Output
	expectedFile := File{
		Action:      "Add",
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
	file1 := File{
		Source:      "test/file",
		Destination: "~/.cache/synchronex/file",
	}
	file2 := File{
		Source:      "test/file2",
		Destination: "~/.cache/synchronex/file2",
	}

	state := &Nex{
		Files: []File{
			file1,
			file2,
		},
	}
	plan := &Nex{
		Files: []File{
			file1,
		},
	}

	// Expected Output
	expectedFile := File{
		Action:      "Remove",
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
	file1 := File{
		Source:      "test/file",
		Destination: "~/.cache/synchronex/file",
	}
	newFile1 := File{
		Source:      "test/file",
		Destination: "~/.cache/synchronex/file",
		User:        "test",
	}

	state := &Nex{
		Files: []File{
			file1,
		},
	}
	plan := &Nex{
		Files: []File{
			newFile1,
		},
	}

	// Expected Output
	expectedFiles := []File{
		File{
			Action:      "Remove",
			Destination: "~/.cache/synchronex/file",
		},
		File{
			Action:      "Add",
			Destination: "~/.cache/synchronex/file",
		},
	}

	// Work
	result := plan.DifferencesFromState(*state)

	// Verification
	assert.True(t, len(result.Files) == 2)
	assert.Equal(t, result.Files, expectedFiles)
}
