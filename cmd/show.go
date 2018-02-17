package cmd

import (
	"strconv"

	"github.com/brunograsselli/lgtm/lgtm"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show pull request on the browser",
	Long:  "Show pull request on the browser",
	Run: func(cmd *cobra.Command, args []string) {
		number, err := strconv.Atoi(args[0])

		if err != nil {
			panic(err)
		}

		browser := lgtm.Browser{
			LastResultsPath: "/tmp/lgtm.json",
		}

		err = browser.Open(int32(number))

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
