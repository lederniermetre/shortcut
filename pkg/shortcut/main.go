package shortcut

import (
	"context"
	"os"
	"time"

	"log/slog"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/lederniermetre/shortcut/pkg/shortcut/gen/client"
	"github.com/lederniermetre/shortcut/pkg/shortcut/gen/client/operations"
	"github.com/lederniermetre/shortcut/pkg/shortcut/gen/models"
	"gitlab.com/greyxor/slogor"
)

const CTX_TIMEOUT = 5000 * time.Millisecond

var clientSC *apiclient.ShortcutAPI

func GetClient() *apiclient.ShortcutAPI {
	if clientSC != nil {
		return clientSC
	}

	// Hack to parse end_date "2023-01-19"
	strfmt.DateTimeFormats = append(strfmt.DateTimeFormats, time.DateOnly)

	// create the transport
	transport := httptransport.New("api.app.shortcut.com", "", nil)
	clientSC = apiclient.New(transport, strfmt.Default)
	return clientSC
}

func GetAuth() runtime.ClientAuthInfoWriter {
	if os.Getenv("SHORTCUT_API_TOKEN") == "" {
		slog.Error("SHORTCUT_API_TOKEN is empty")
		os.Exit(1)
	}

	return httptransport.APIKeyAuth("Shortcut-Token", "header", os.Getenv("SHORTCUT_API_TOKEN"))
}

func RetrieveIteration(name string) models.IterationSlim {
	searchIterationsParams := &operations.SearchIterationsParams{}
	search := &models.Search{
		Detail:   "slim",
		Query:    &name,
		PageSize: 1,
	}

	err := search.Validate(strfmt.Default)
	if err != nil {
		slog.Error("Search is invalid", slogor.Err(err))
		os.Exit(1)
	}

	searchIterationsParams.Search = search

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	searchIterationsParams.SetContext(ctx)

	searchResult, err := GetClient().Operations.SearchIterations(searchIterationsParams, GetAuth())
	if err != nil {
		slog.Error("Can not retrieve search", slogor.Err(err), slog.String("name", name))
		os.Exit(1)
	}

	if len(searchResult.Payload.Data) < 1 {
		slog.Error("Search has retrieve no result", slog.String("name", name))
		os.Exit(1)
	}

	return *searchResult.Payload.Data[0]
}

func StoriesByIteration(iterationID int64) []*models.StorySlim {
	listIterationStoriesParams := &operations.ListIterationStoriesParams{
		IterationPublicID: iterationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	listIterationStoriesParams.SetContext(ctx)

	allStories, err := GetClient().Operations.ListIterationStories(listIterationStoriesParams, GetAuth())
	if err != nil {
		slog.Error("Can not retrieve iteration", slogor.Err(err))
		os.Exit(1)
	}

	return allStories.Payload
}

func GetMember(uuid strfmt.UUID) models.Member {
	getMemberParams := &operations.GetMemberParams{
		MemberPublicID: uuid,
	}

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	getMemberParams.SetContext(ctx)

	ownerInfo, err := GetClient().Operations.GetMember(getMemberParams, GetAuth())
	if err != nil {
		slog.Error("can not retrieve iteration", "detail", err.Error())
		os.Exit(1)
	}

	return *ownerInfo.Payload
}

func GetWorkflow(id int64) models.Workflow {
	getWorkflowParams := &operations.GetWorkflowParams{
		WorkflowPublicID: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	getWorkflowParams.SetContext(ctx)

	workflow, err := GetClient().Operations.GetWorkflow(getWorkflowParams, GetAuth())
	if err != nil {
		slog.Error("Can not retrieve workflow", slogor.Err(err))
		os.Exit(1)
	}

	return *workflow.Payload
}

func GetEpic(id int64) models.Epic {
	getEpicParams := &operations.GetEpicParams{
		EpicPublicID: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	getEpicParams.SetContext(ctx)

	epic, err := GetClient().Operations.GetEpic(getEpicParams, GetAuth())
	if err != nil {
		slog.Error("Can not retrieve epic", slogor.Err(err))
		os.Exit(1)
	}

	return *epic.Payload
}
