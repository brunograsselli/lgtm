package cmd

import (
	"fmt"

	"github.com/brunograsselli/lgtm"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests waiting for your review",
	Long:  `List pull requests waiting for your review.`,
	Run: func(cmd *cobra.Command, args []string) {
		showAll, err := cmd.Flags().GetBool("all")

		if err != nil {
			panic(err)
		}

		err = lgtm.List(showAll)

		if err != nil {
			fmt.Printf("An error occurred while processing your request:\n  %s\n", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("all", "a", false, "List all open pull requests")
}
