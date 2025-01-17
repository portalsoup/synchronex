package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"synchronex/common"
)

var PlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Personal computer state manager",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		plan := *Nex
		state, err := common.ReadStatefile()
		if err != nil {
			log.Fatal(err)
		}

		diff := common.PrintPretty(plan.DifferencesFromState(*state))
		fmt.Println(diff)
	},
}
