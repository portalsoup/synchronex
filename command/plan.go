package command

import (
	"github.com/spf13/cobra"
)

var PlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Personal computer state manager",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
