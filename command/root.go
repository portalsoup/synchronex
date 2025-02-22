package command

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"synchronex/common"
	"synchronex/schema"
)

var (
	Nex     *schema.Nex
	State   *schema.Nex
	Diff    *schema.Nex
	RootCmd = &cobra.Command{
		Use:   "synchronex",
		Short: "Personal computer state manager",
		Long:  ``,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			_, err := common.ReadStatefile()
			if err != nil {
				// Initialize empty state if necessary
				err = common.WriteStatefile(schema.Nex{
					Files: []schema.File{},
				})
				if err != nil {
					log.Fatalf("Error initializing state file: %s", err)
				}
			}

			quiet, err := cmd.Root().Flags().GetBool("quiet")
			if err != nil {
				log.Fatalf("Error getting verbose flag: %s", err)
			}

			common.ConfigureLogger(!quiet)

			if len(args) == 0 {
				log.Fatal("No nex file specified")
			}

			// Check if the file exists
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				log.Fatalf("File does not exist: %s", args[0])
			}

			// Convert to absolute path
			absPath, err := filepath.Abs(args[0])
			if err != nil {
				log.Fatalf("Error converting path to absolute path: %s", err)
			}

			parsedNex, err := common.ParseNexFile(absPath)
			if err != nil || parsedNex == nil {
				log.Fatalf("Error parsing nex file[%s]:\n\tReason:\n%s", parsedNex, err)
			}
			parsedNex.ExpandHomeFolder()
			Nex = parsedNex

			// Validate changes with user
			plan := parsedNex
			State, err = common.ReadStatefile()
			if err != nil {
				log.Fatal(err)
			}

			diff := plan.DifferencesFromState(*State)
			Diff = diff
		},
	}
)

func init() {
	RootCmd.Flags().BoolP("quiet", "q", false, "Silence output to stdout")

	RootCmd.AddCommand(PlanCmd)
	RootCmd.AddCommand(ApplyCmd)
}

func Execute() (err error) {
	err = RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
	return err
}
