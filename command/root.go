package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"synchronex/common"
	"synchronex/schema"
)

var (
	nexAbsPath string
	nexes      []*schema.Nex
	rootCmd    = &cobra.Command{
		Use:   "synchronex",
		Short: "Personal computer state manager",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			for _, nex := range nexes {
				log.Printf("Got this nex: %#v", *nex)
			}
		},
	}
)

func Execute() error {

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			dirPath, err := os.Getwd()
			if err != nil {
				return err
			}

			foundNexes, err := common.GetNexes(dirPath)
			if err != nil {
				return fmt.Errorf("Failed to load Nex from path %s. %v", foundNexes, err)
			}
			nexes = foundNexes
		}

		for _, nexPath := range args {
			foundNexes, err := common.GetNexes(nexPath)
			if err != nil {
				return fmt.Errorf("Failed to load Nex from path %s. %v", nexPath, err)
			}
			nexes = append(nexes, foundNexes...) // error foundNexes is an array
		}
		return nil
	}

	return rootCmd.Execute()
}
