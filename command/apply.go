package command

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"synchronex/common"
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

		diff := common.PrintPretty(plan.DifferencesFromState(*state))
		if diff == "{}" {
			log.Println("No changes to apply.")
			os.Exit(0)
		}
		userAccept := promptUser(fmt.Sprintf("%s\nThe following changes will be made, Continue? (y/n)", diff), "y")
		if !userAccept {
			log.Println("User denied the changes...")
			os.Exit(0)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Doing the work!")
		common.ApplyDiff(Diff.Files)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		nex := *Nex
		err := common.WriteStatefile(nex)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func promptUser(message string, token string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message)
	input, _ := reader.ReadString('\n')

	if strings.TrimSpace(input) == token {
		fmt.Println("Continuing with the changes...")
		return true
		// Proceed with the rest of the code
	} else {
		fmt.Println("Operation aborted by the user.")
		return false
	}
}
