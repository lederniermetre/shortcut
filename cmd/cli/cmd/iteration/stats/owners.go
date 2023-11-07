package stats

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-openapi/strfmt"
	"github.com/lederniermetre/shortcut/pkg/shortcut"
	"github.com/spf13/cobra"
)

var ownersCmd = &cobra.Command{
	Use:   "owners",
	Short: "Compute stats based on stories for owners",
	Long: `This tasks will base computation on stories in retrieve iteration.

The load is the sum of each story assign to an owner.
Estimate is divided by number of owners when multi-tenancy`,
	Run: func(cmd *cobra.Command, args []string) {
		iterationFlag := cmd.Parent().Parent().PersistentFlags().Lookup("iteration")
		if iterationFlag == nil {
			slog.Error("Can not retrieved iteration flag")
			os.Exit(1)
		}

		iterationName := iterationFlag.Value.String()
		slog.Debug("Working on iteration", slog.String("name", iterationName))

		iteration := shortcut.RetrieveIteration(iterationName)

		slog.Info("Iteration retrieved", slog.String("name", *iteration.Name))

		allStories := shortcut.StoriesByIteration(*iteration.ID)

		ownersUUID := map[strfmt.UUID]int64{}
		for _, story := range allStories {
			if story.Estimate == nil {
				slog.Warn("Story assign but not estimated", slog.String("name", *story.Name))
				continue
			}

			slog.Debug(
				"Compute story",
				slog.String("name", *story.Name),
				slog.Int("owners", len(story.OwnerIds)),
				slog.Int64("estimate", *story.Estimate),
			)

			if len(story.OwnerIds) == 0 {
				slog.Warn("Story has no owners", slog.String("name", *story.Name))
				continue
			}

			estimate := *story.Estimate
			if len(story.OwnerIds) > 1 {
				slog.Debug("Story shared, split estimate", slog.String("name", *story.Name))
				estimate = estimate / int64(len(story.OwnerIds))
			}

			for _, ownedId := range story.OwnerIds {
				val, init := ownersUUID[ownedId]
				if init {
					ownersUUID[ownedId] = val + estimate
				} else {
					ownersUUID[ownedId] = estimate
				}
			}
		}

		slog.Info("===== Load by owners =====")

		for ownerUUID, load := range ownersUUID {
			slog.Info(fmt.Sprintf("%s has %d of load", *shortcut.GetMember(ownerUUID).Profile.Name, load))
		}
	},
}

func newOwnersCommand() *cobra.Command {
	return ownersCmd
}

func init() {
}
