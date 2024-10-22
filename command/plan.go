package command

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	planCmd = &cobra.Command{
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
