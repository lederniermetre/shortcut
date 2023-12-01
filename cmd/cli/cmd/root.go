/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/lederniermetre/shortcut/cmd/cli/cmd/iteration"
	"github.com/lederniermetre/shortcut/pkg/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shortcut",
	Short: "Supercharge your shortcut",
	Long:  `This cli aims to be your best buddy and get a lots of stats`,
}
var debugCmd bool

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initCmd)
	rootCmd.PersistentFlags().BoolVarP(&debugCmd, "debug", "d", false, "Active debug log level")

	rootCmd.AddCommand(iteration.NewCommand())
}

func initCmd() {
	utils.SetLogger(debugCmd)
}
