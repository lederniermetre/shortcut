package stats

import (
	"log/slog"

	"github.com/spf13/cobra"
	"gitlab.com/greyxor/slogor"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Compute statistics for dedicated iteration",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Usage()
		if err != nil {
			slog.Error("Can not display usage", slogor.Err(err))
		}
	},
}

func NewCommand() *cobra.Command {
	return statsCmd
}

func init() {
	statsCmd.AddCommand(newOwnersCommand())
	statsCmd.AddCommand(newStoriesCommand())
	statsCmd.PersistentFlags().StringP("query", "q", "Ops", "Search query")
}
