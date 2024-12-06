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
	RootCmd = &cobra.Command{
		Use:   "synchronex",
		Short: "Personal computer state manager",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Running command: %s", cmd.Name())
			log.Printf("Got this nex: %#v", *Nex)
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
			Nex = parsedNex
		},
	}
)

func init() {
	RootCmd.Flags().BoolP("quiet", "q", false, "Silence output to stdout")
}

func Execute() (err error) {
	err = RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
	return err
}
