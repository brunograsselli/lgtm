package cmd

import (
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

		lgtm.List(showAll)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("all", "a", false, "List all open pull requests")
}
