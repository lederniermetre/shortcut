package stats

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/lederniermetre/shortcut/pkg/shortcut"
	"github.com/lederniermetre/shortcut/pkg/shortcut/gen/client/operations"
	"github.com/lederniermetre/shortcut/pkg/shortcut/gen/models"
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

		clientSC := shortcut.GetClient()
		apiKeyHeaderAuth := shortcut.GetAuth()

		searchIterationsParams := &operations.SearchIterationsParams{}
		search := &models.Search{
			Detail:   "slim",
			Query:    &iterationName,
			PageSize: 1,
		}

		err := search.Validate(strfmt.Default)
		if err != nil {
			slog.Error("Search is invalid", slogor.Err(err))
			os.Exit(1)
		}

		searchIterationsParams.Search = search

		ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
		defer cancel()

		searchIterationsParams.SetContext(ctx)

		searchResult, err := clientSC.Operations.SearchIterations(searchIterationsParams, apiKeyHeaderAuth)
		if err != nil {
			slog.Error("Can not retrieve search", slogor.Err(err), slog.String("name", iterationName))
			os.Exit(1)
		}

		if len(searchResult.Payload.Data) < 1 {
			slog.Error("Search has retrieve no result", slog.String("name", iterationName))
			os.Exit(1)
		}

		slog.Info("Retrieve iteration informations", slog.String("name", *searchResult.Payload.Data[0].Name))

		listIterationStoriesParams := &operations.ListIterationStoriesParams{
			IterationPublicID: *searchResult.Payload.Data[0].ID,
		}
		listIterationStoriesParams.SetContext(ctx)

		allStories, err := clientSC.Operations.ListIterationStories(listIterationStoriesParams, apiKeyHeaderAuth)
		if err != nil {
			slog.Error("Can not retrieve iteration", slogor.Err(err))
			os.Exit(1)
		}

		ownersUUID := map[strfmt.UUID]int64{}
		for _, story := range allStories.Payload {
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
			getMemberParams := &operations.GetMemberParams{
				MemberPublicID: ownerUUID,
			}
			getMemberParams.SetContext(ctx)

			ownerInfo, err := clientSC.Operations.GetMember(getMemberParams, apiKeyHeaderAuth)
			if err != nil {
				slog.Error("can not retrieve iteration", "detail", err.Error())
				os.Exit(1)
			}

			slog.Info(fmt.Sprintf("%s has %d of load", *ownerInfo.Payload.Profile.Name, load))
		}
	},
}

func newOwnersCommand() *cobra.Command {
	return ownersCmd
}

func init() {
}
