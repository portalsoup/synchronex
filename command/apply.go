package command

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"synchronex/common"
	"synchronex/schema"
)

var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Personal computer state manager",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Validate changes with user
		plan := *Nex
		state, err := common.ReadStatefile()
		if err != nil {
			log.Fatal(err)
		}

		rawDiff := plan.DifferencesFromState(*state)
		add, remove := rawDiff.DiffSummary()
		var summary string
		if add != 0 || remove != 0 {
			summary = common.PrintPretty(plan.DifferencesFromState(*state))
			fmt.Println(fmt.Sprintf("%s\nThese changes will be made\n%d to add, %d to remove\nContinue? (y/n)", summary, add, remove))
			userAccept := promptUser("y")
			if !userAccept {
				fmt.Println("User denied the changes...")
				os.Exit(0)
			}
		} else {
			fmt.Println("No changes to apply.")
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ApplyDiff(Diff.Files)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		nex := *Nex
		err := common.WriteStatefile(nex)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func promptUser(token string) bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	if strings.TrimSpace(input) == token {
		return true
		// Proceed with the rest of the code
	} else {
		fmt.Println("Operation aborted by the user.")
		return false
	}
}

// applyDiff applies the given list of file changes (diff).
func ApplyDiff(diff []schema.File) error {
	fmt.Printf("Applying %d changes...\n", len(diff))
	for _, f := range diff {
		var err error
		switch f.Action {
		case "Add":
			err = applyAdd(f)
		case "Remove":
			err = applyRemove(f)
		default:
			log.Println("Unknown action:", f.Action)
			continue
		}

		// Log error if an operation fails
		if err != nil {
			log.Printf("Failed to apply action %s for %s: %v\n", f.Action, f.Destination, err)
			return err
		}
	}

	return nil
}

// applyAdd creates a file or folder at the given destination.
// If a source is provided, it copies the source content to the destination.
func applyAdd(f schema.File) error {
	if f.Source != "" {
		// If source is specified, copy the file to the destination
		srcFile, err := os.Open(f.Source)
		if err != nil {
			if os.IsNotExist(err) {
				srcFile, err = os.Create(f.Source)
				if err == nil {
					return err
				}
			} else {
				return err
			}
		}
		defer srcFile.Close()

		// Ensure target directory exists
		destDir := path.Dir(f.Destination)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return err
		}

		// Create the destination file
		destFile, err := os.Create(f.Destination)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Copy the source file content to the destination
		if _, err := io.Copy(destFile, srcFile); err != nil {
			return err
		}
	} else {
		// If source is not specified, create an empty folder
		if err := os.MkdirAll(f.Destination, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// applyRemove deletes the specified file or folder.
func applyRemove(f schema.File) error {
	info, err := os.Stat(f.Destination)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // The file/folder doesn't exist, nothing to do
		}
		return err
	}

	if info.IsDir() {
		// Remove the folder and its contents
		if err := os.RemoveAll(f.Destination); err != nil {
			return err
		}
	} else {
		// Remove the file
		if err := os.Remove(f.Destination); err != nil {
			return err
		}
	}

	return nil
}
