package cmd

import (
	"fmt"
	"sort"

	"github.com/brunograsselli/lgtm/lgtm"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration",
	Long:  "Show configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config := &lgtm.Config{}
		secrets := &lgtm.Secrets{Path: secretsPath}

		showConfig(config, secrets)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func showConfig(config *lgtm.Config, secrets *lgtm.Secrets) {
	repos := config.Repos()
	sort.Strings(repos)
	token, _ := secrets.Token()

	fmt.Printf("User: %s\n", config.UserName())

	if token != nil {
		fmt.Println("Logged In: yes")
	} else {
		fmt.Println("Logged In: no")
	}

	fmt.Println("Repositories:")

	for _, repo := range repos {
		fmt.Printf("  - %s\n", repo)
	}
}
