package cmd

import (
	"fmt"
	"sort"

	"github.com/brunograsselli/lgtm/lgtm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration",
	Long:  "Show configuration",
	Run: func(cmd *cobra.Command, args []string) {
		repos := &lgtm.Repos{}
		secrets := &lgtm.Secrets{Path: secretsPath}

		showConfig(repos, secrets)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func showConfig(repos *lgtm.Repos, secrets *lgtm.Secrets) {
	all := repos.All()
	sort.Strings(all)
	token, _ := secrets.Token()

	fmt.Printf("User: %s\n", viper.GetString("username"))

	if token != nil {
		fmt.Println("Logged In: yes")
	} else {
		fmt.Println("Logged In: no")
	}

	fmt.Println("Repositories:")

	for _, repo := range all {
		fmt.Printf("  - %s\n", repo)
	}
}
