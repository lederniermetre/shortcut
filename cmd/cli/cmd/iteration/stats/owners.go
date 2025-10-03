package stats

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-openapi/strfmt"
	"github.com/lederniermetre/shortcut/pkg/shortcut"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var ownersCmd = &cobra.Command{
	Use:   "owners",
	Short: "Compute stats based on stories for owners",
	Long: `This tasks will base computation on stories in retrieve iteration.

The load is the sum of each story assign to an owner.
Estimate is divided by number of owners when multi-tenancy`,
	Run: func(cmd *cobra.Command, args []string) {
		queryFlag := cmd.Parent().PersistentFlags().Lookup("query")
		if queryFlag == nil {
			slog.Error("Can not retrieved iteration flag")
			os.Exit(1)
		}
		limitFlag, err := cmd.Parent().Flags().GetInt("limit")
		if err != nil {
			slog.Error("Can not retrieved limit flag", slog.Any("error", err))
			os.Exit(1)
		}

		shortcutQuery := queryFlag.Value.String()
		slog.Debug("Search", slog.String("name", shortcutQuery))

		iterations, err := shortcut.RetrieveIterations(shortcutQuery, limitFlag, "")
		if err != nil {
			slog.Error("Cannot retrieve iterations", slog.Any("error", err))
			os.Exit(1)
		}

		ownersUUID := map[strfmt.UUID]int64{}
		for _, iteration := range iterations {
			slog.Info("Iteration retrieved", slog.String("name", *iteration.Name))

			allStories, err := shortcut.StoriesByIteration(*iteration.ID)
			if err != nil {
				slog.Error("Cannot retrieve stories", slog.Any("error", err))
				os.Exit(1)
			}

			for _, story := range allStories {
				if story.Archived != nil && *story.Archived {
					pterm.Info.Printfln("Story %s is archived skipping", *story.Name)
					continue
				}

				if story.Estimate == nil {
					pterm.Warning.Printfln("Story assign but not estimated: %s", *story.Name)
					continue
				}

				slog.Debug(
					"Compute story",
					slog.String("name", *story.Name),
					slog.Int("owners", len(story.OwnerIds)),
					slog.Int64("estimate", *story.Estimate),
				)

				if len(story.OwnerIds) == 0 {
					pterm.Warning.Printfln("Story has no owners: %s", *story.Name)
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
		}

		pterm.DefaultHeader.WithFullWidth().Println("Load by owners")

		ordererOwnersUUID := shortcut.OrdererOwnersUUID(ownersUUID)
		var ptermBar []pterm.Bar

		for _, owner := range ordererOwnersUUID {
			member, err := shortcut.GetMember(owner.UUID)
			if err != nil {
				slog.Error("Cannot retrieve member", slog.Any("error", err))
				os.Exit(1)
			}
			memberName := *member.Profile.Name
			slog.Debug(fmt.Sprintf("%s has %d of load", memberName, owner.Load))
			ptermBar = append(ptermBar, pterm.Bar{Label: memberName, Value: int(owner.Load)})
		}

		err = pterm.DefaultBarChart.WithHorizontal().WithBars(ptermBar).WithWidth(15).WithShowValue().Render()
		if err != nil {
			slog.Error("Rendering epics table", slog.Any("error", err))
		}
	},
}

func newOwnersCommand() *cobra.Command {
	return ownersCmd
}

func init() {
}
