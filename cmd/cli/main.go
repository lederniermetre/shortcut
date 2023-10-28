package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/lederniermetre/shortcut/pkg/shortcut/gen/client"
	"github.com/lederniermetre/shortcut/pkg/shortcut/gen/client/operations"
	"github.com/lederniermetre/shortcut/pkg/shortcut/gen/models"

	"gitlab.com/greyxor/slogor"
)

func main() {
	iterationName := flag.String("iteration", "Ops", "Iteration title you are looking for")
	debug := flag.Bool("debug", false, "Display debug logs")
	flag.Parse()

	logLevel := slog.LevelInfo
	if *debug {
		logLevel = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slogor.NewHandler(os.Stderr, &slogor.Options{
		TimeFormat: time.Stamp,
		Level:      logLevel,
		ShowSource: true,
	})))

	// create the transport
	transport := httptransport.New("api.app.shortcut.com", "", nil)
	clientSC := apiclient.New(transport, strfmt.Default)
	apiKeyHeaderAuth := httptransport.APIKeyAuth("Shortcut-Token", "header", os.Getenv("SHORTCUT_API_TOKEN"))

	searchIterationsParams := &operations.SearchIterationsParams{}
	search := &models.Search{
		Detail:   "slim",
		Query:    iterationName,
		PageSize: 1,
	}

	err := search.Validate(strfmt.Default)
	if err != nil {
		slog.Error("Search is invalid", slogor.Err(err))
		os.Exit(1)
	}

	searchIterationsParams.Search = search

	// Hack to parse end_date "2023-01-19"
	strfmt.DateTimeFormats = append(strfmt.DateTimeFormats, time.DateOnly)

	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	searchIterationsParams.SetContext(ctx)

	searchResult, err := clientSC.Operations.SearchIterations(searchIterationsParams, apiKeyHeaderAuth)
	if err != nil {
		slog.Error("Can not retrieve search", slogor.Err(err), slog.String("name", *iterationName))
		os.Exit(1)
	}

	if len(searchResult.Payload.Data) < 1 {
		slog.Error("Search has retrieve no result", slog.String("name", *iterationName))
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
	for _, it := range allStories.Payload {
		if it.Estimate == nil {
			slog.Warn(fmt.Sprintf("OMG no estimate on story: %s", *it.Name))
			continue
		}

		slog.Debug(
			"Compute story",
			slog.String("name", *it.Name),
			slog.Int("owners", len(it.OwnerIds)),
			slog.Int64("estimate", *it.Estimate),
		)

		if len(it.OwnerIds) == 0 {
			slog.Warn(fmt.Sprintf("Story has no owners"), slog.String("name", *it.Name))
			continue
		}

		estimate := *it.Estimate
		if len(it.OwnerIds) > 1 {
			slog.Debug("Story shared, split estimate", slog.String("name", *it.Name))
			estimate = estimate / 2
		}

		for _, ownedId := range it.OwnerIds {
			val, init := ownersUUID[ownedId]
			if init {
				ownersUUID[ownedId] = val + estimate
			} else {
				ownersUUID[ownedId] = estimate
			}
		}
	}

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
}
