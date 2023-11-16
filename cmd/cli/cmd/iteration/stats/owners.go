package stats

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-openapi/strfmt"
	"github.com/lederniermetre/shortcut/pkg/shortcut"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gitlab.com/greyxor/slogor"
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

		pterm.DefaultHeader.WithFullWidth().Println("Load by owners")

		var ptermBar []pterm.Bar

		for ownerUUID, load := range ownersUUID {
			memberName := *shortcut.GetMember(ownerUUID).Profile.Name
			slog.Debug(fmt.Sprintf("%s has %d of load", memberName, load))
			ptermBar = append(ptermBar, pterm.Bar{Label: memberName, Value: int(load)})
		}

		err := pterm.DefaultBarChart.WithHorizontal().WithBars(ptermBar).WithWidth(15).WithShowValue().Render()
		if err != nil {
			slog.Error("Rendering epics table", slogor.Err(err))
		}
	},
}

func newOwnersCommand() *cobra.Command {
	return ownersCmd
}

func init() {
}
