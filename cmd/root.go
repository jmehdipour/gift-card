package cmd

import (
	"github.com/jmehdipour/gift-card/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Flag variables
	cfgFile string

	// rootCMD represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:              "gift-card",
		Short:            "gift-card service",
		PersistentPreRun: preRun,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(databaseCMD)
}

func preRun(_ *cobra.Command, _ []string) {
	config.Init(cfgFile)
}

func Execute() error {
	return rootCmd.Execute()
}
