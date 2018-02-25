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
		secrets := &lgtm.Secrets{Path: secretsPath}
		config := lgtm.NewConfig()

		err := lgtm.Login(secrets, config)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
