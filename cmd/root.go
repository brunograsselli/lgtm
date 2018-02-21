package cmd

import (
	"fmt"
	"os"

	"github.com/brunograsselli/lgtm/lgtm"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "lgtm",
	Short: "Watch pull requests waiting for your review",
	Long:  "Watch pull requests waiting for your review",
	Run: func(cmd *cobra.Command, args []string) {
		showAll, err := cmd.Flags().GetBool("all")

		if err != nil {
			panic(err)
		}

		secrets := &lgtm.Secrets{Path: secretsPath}
		config := lgtm.NewConfig()

		err = lgtm.List(showAll, secrets, config)

		if err != nil {
			fmt.Printf("An error occurred while processing your request:\n  %s\n", err.Error())
		}
	},
}

var secretsPath = fmt.Sprintf("%s/.lgtm.secret", os.Getenv("HOME"))

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lgtm.yml)")
	rootCmd.Flags().BoolP("all", "a", false, "List all open pull requests")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".lgtm" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lgtm")
	}

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
