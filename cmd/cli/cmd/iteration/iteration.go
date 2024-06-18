/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package iteration

import (
	"log/slog"

	"github.com/lederniermetre/shortcut/cmd/cli/cmd/iteration/stats"
	"github.com/spf13/cobra"
	"gitlab.com/greyxor/slogor"
)

// iterationCmd represents the iteration command
var iterationCmd = &cobra.Command{
	Use:   "iteration",
	Short: "Work on iteration entities",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Usage()
		if err != nil {
			slog.Error("Can not display usage", slogor.Err(err))
		}
	},
}

func NewCommand() *cobra.Command {
	return iterationCmd
}

func init() {
	iterationCmd.AddCommand(stats.NewCommand())
	iterationCmd.PersistentFlags().StringP("query", "q", "Ops", "Search query")
}
