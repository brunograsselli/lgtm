package cmd

import (
	"github.com/brunograsselli/lgtm/lgtm"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout",
	Long:  "Logout",
	Run: func(cmd *cobra.Command, args []string) {
		secrets := lgtm.NewSecrets(secretsPath)

		if err := secrets.DeleteToken(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
