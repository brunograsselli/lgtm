package cmd

import (
	"github.com/brunograsselli/lgtm"
	"github.com/spf13/cobra"
)

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "List watched repositories",
	Long:  "List watched repositories",
	Run: func(cmd *cobra.Command, args []string) {
		repos := &lgtm.Repos{}

		repos.List()
	},
}

func init() {
	rootCmd.AddCommand(reposCmd)
}
