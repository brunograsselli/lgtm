package cmd

import (
	"fmt"
	"os"

	"github.com/brunograsselli/lgtm/lgtm"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login",
	Long:  "Login",
	Run: func(cmd *cobra.Command, args []string) {
		secrets := lgtm.NewSecrets(secretsPath)
		config := lgtm.NewConfig()

		if err := lgtm.Login(secrets, config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
