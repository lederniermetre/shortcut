package stats

import (
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Compute statistics for dedicated iteration",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func NewCommand() *cobra.Command {
	return statsCmd
}

func init() {
	statsCmd.AddCommand(newOwnersCommand())
	statsCmd.AddCommand(newStoriesCommand())
}
