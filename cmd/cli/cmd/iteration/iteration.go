/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package iteration

import (
	"github.com/lederniermetre/shortcut/cmd/cli/cmd/iteration/stats"
	"github.com/spf13/cobra"
)

// iterationCmd represents the iteration command
var iterationCmd = &cobra.Command{
	Use:   "iteration",
	Short: "Work on iteration entities",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func NewCommand() *cobra.Command {
	return iterationCmd
}

func init() {
	iterationCmd.AddCommand(stats.NewCommand())
	iterationCmd.PersistentFlags().StringP("iteration", "i", "Ops", "Iteration title you are looking for")
}
