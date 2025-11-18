package shortcut

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/lederniermetre/shortcut/pkg/shortcut/gen/client"
	"github.com/lederniermetre/shortcut/pkg/shortcut/gen/client/operations"
	"github.com/lederniermetre/shortcut/pkg/shortcut/gen/models"
)

const CTX_TIMEOUT = 5000 * time.Millisecond

var (
	clientSC   *apiclient.ShortcutAPI
	clientOnce sync.Once
)

func GetClient() *apiclient.ShortcutAPI {
	clientOnce.Do(func() {
		// Hack to parse end_date "2023-01-19"
		strfmt.DateTimeFormats = append(strfmt.DateTimeFormats, time.DateOnly)

		// create the transport
		transport := httptransport.New("api.app.shortcut.com", "", nil)
		clientSC = apiclient.New(transport, strfmt.Default)
	})

	return clientSC
}

func GetAuth() (runtime.ClientAuthInfoWriter, error) {
	token := os.Getenv("SHORTCUT_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("SHORTCUT_API_TOKEN environment variable is required")
	}

	return httptransport.APIKeyAuth("Shortcut-Token", "header", token), nil
}

func getIterationData(params *operations.SearchIterationsParams, query string) ([]*models.IterationSlim, error) {
	var data []*models.IterationSlim

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	params.SetContext(ctx)

	auth, err := GetAuth()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	searchResult, err := GetClient().Operations.SearchIterations(params, auth)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve search for query %q: %w", query, err)
	}

	if len(searchResult.Payload.Data) < 1 {
		return nil, fmt.Errorf("search returned no results for query %q", query)
	}

	data = searchResult.Payload.Data

	if *searchResult.GetPayload().Total > *params.PageSize && searchResult.GetPayload().Next != nil {
		slog.Debug("Continu on next page", slog.String("url", *searchResult.GetPayload().Next))
		nextData, err := RetrieveIterations(query, int(*params.PageSize), *searchResult.GetPayload().Next)
		if err != nil {
			return nil, err
		}
		data = append(data, nextData...)
	}

	return data, nil
}

func RetrieveIterations(query string, pageLimit int, nextURL string) ([]*models.IterationSlim, error) {
	searchDetail := "slim"
	pageSize := int64(pageLimit)
	var nextToken string

	if nextURL != "" {
		u, err := url.Parse(nextURL)
		if err != nil {
			return nil, fmt.Errorf("cannot parse URL %q: %w", nextURL, err)
		}

		nextToken = u.Query().Get("next")
	}

	searchIterationsParams := &operations.SearchIterationsParams{
		Detail:   &searchDetail,
		Query:    query,
		PageSize: &pageSize,
		Next:     &nextToken,
	}

	data, err := getIterationData(searchIterationsParams, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func StoriesByIteration(iterationID int64) ([]*models.StorySlim, error) {
	listIterationStoriesParams := &operations.ListIterationStoriesParams{
		IterationPublicID: iterationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	listIterationStoriesParams.SetContext(ctx)

	auth, err := GetAuth()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	allStories, err := GetClient().Operations.ListIterationStories(listIterationStoriesParams, auth)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve stories for iteration %d: %w", iterationID, err)
	}

	return allStories.Payload, nil
}

func GetMember(uuid strfmt.UUID) (*models.Member, error) {
	getMemberParams := &operations.GetMemberParams{
		MemberPublicID: uuid,
	}

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	getMemberParams.SetContext(ctx)

	auth, err := GetAuth()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	ownerInfo, err := GetClient().Operations.GetMember(getMemberParams, auth)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve member %s: %w", uuid, err)
	}

	return ownerInfo.Payload, nil
}

func GetWorkflow(id int64) (*models.Workflow, error) {
	getWorkflowParams := &operations.GetWorkflowParams{
		WorkflowPublicID: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	getWorkflowParams.SetContext(ctx)

	auth, err := GetAuth()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	workflow, err := GetClient().Operations.GetWorkflow(getWorkflowParams, auth)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve workflow %d: %w", id, err)
	}

	return workflow.Payload, nil
}

func GetEpic(id int64) (*models.Epic, error) {
	getEpicParams := &operations.GetEpicParams{
		EpicPublicID: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), CTX_TIMEOUT)
	defer cancel()
	getEpicParams.SetContext(ctx)

	auth, err := GetAuth()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	epic, err := GetClient().Operations.GetEpic(getEpicParams, auth)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve epic %d: %w", id, err)
	}

	return epic.Payload, nil
}

func IncreaseEpicsStoriesCounter(storyWorkflowState WorkflowInfo, epicsStats EpicsStats) EpicsStats {
	if storyWorkflowState.Type == "started" {
		epicsStats.StoriesStarted++
		return epicsStats
	}

	if storyWorkflowState.Type == "unstarted" {
		epicsStats.StoriesUnstarted++
		return epicsStats
	}

	if storyWorkflowState.Type == "done" {
		epicsStats.StoriesDone++
		return epicsStats
	}

	if storyWorkflowState.Type == "backlog" {
		epicsStats.StoriesBacklog++
		slog.Warn("Story is in backlog state this is not normal")
		return epicsStats
	}

	slog.Error("Worfklow state type unknown", slog.String("name", storyWorkflowState.Type), slog.String("type", storyWorkflowState.Name))

	return epicsStats
}

func IncreaseEpicsEstimateCounter(storyWorkflowState WorkflowInfo, epicsStats EpicsStats, estimate int) EpicsStats {
	if storyWorkflowState.Type == "started" {
		epicsStats.EstimateStarted += estimate
		return epicsStats
	}

	if storyWorkflowState.Type == "unstarted" {
		epicsStats.EstimateUnstarted += estimate
		return epicsStats
	}

	if storyWorkflowState.Type == "done" {
		epicsStats.EstimateDone += estimate
		return epicsStats
	}

	if storyWorkflowState.Type == "backlog" {
		epicsStats.EstimateBacklog += estimate
		slog.Warn("Story is in backlog state this is not normal")
		return epicsStats
	}

	slog.Error("Worfklow state type unknown", slog.String("name", storyWorkflowState.Type), slog.String("type", storyWorkflowState.Name))

	return epicsStats
}

func SummaryEpicStat(epic EpicsStats) EpicsStats {
	totalEpicsStories := epic.StoriesUnstarted + epic.StoriesStarted + epic.StoriesDone + epic.StoriesBacklog
	if totalEpicsStories != 0 {
		epic.StoriesDonePercent = epic.StoriesDone * 100 / totalEpicsStories
		epic.StoriesUnstartedPercent = epic.StoriesUnstarted * 100 / totalEpicsStories
		epic.StoriesStartedPercent = epic.StoriesStarted * 100 / totalEpicsStories
		epic.StoriesBacklogPercent = epic.StoriesBacklog * 100 / totalEpicsStories
	}

	totalEpicsEstimateStories := epic.EstimateUnstarted + epic.EstimateStarted + epic.EstimateDone + epic.EstimateBacklog
	if totalEpicsEstimateStories != 0 {
		epic.EstimateDonePercent = epic.EstimateDone * 100 / totalEpicsEstimateStories
		epic.EstimateUnstartedPercent = epic.EstimateUnstarted * 100 / totalEpicsEstimateStories
		epic.EstimateStartedPercent = epic.EstimateStarted * 100 / totalEpicsEstimateStories
		epic.EstimateBacklogPercent = epic.EstimateBacklog * 100 / totalEpicsEstimateStories
	}

	return epic
}

func OrdererOwnersUUID(ownersUUID map[strfmt.UUID]int64) []OwnerStats {
	var ordererOwnersUUID []OwnerStats

	for k, v := range ownersUUID {
		ordererOwnersUUID = append(ordererOwnersUUID, struct {
			UUID strfmt.UUID
			Load int64
		}{k, v})
	}

	// Sort the slice by values in descending order
	sort.Slice(ordererOwnersUUID, func(i, j int) bool {
		return ordererOwnersUUID[i].Load > ordererOwnersUUID[j].Load
	})

	return ordererOwnersUUID
}

func ComputeEpicGlobalStat(global GlobalEpicStats, epic EpicsStats) GlobalEpicStats {
	global.EstimateUnstarted += epic.EstimateUnstarted
	global.EstimateStarted += epic.EstimateStarted
	global.EstimateDone += epic.EstimateDone
	global.EstimateBacklog += epic.EstimateBacklog

	global.StoriesUnstarted += epic.StoriesUnstarted
	global.StoriesStarted += epic.StoriesStarted
	global.StoriesDone += epic.StoriesDone
	global.StoriesBacklog += epic.StoriesBacklog

	totalEpicsStories := global.StoriesUnstarted + global.StoriesStarted + global.StoriesDone + global.StoriesBacklog
	if totalEpicsStories != 0 {
		global.StoriesDonePercent = global.StoriesDone * 100 / totalEpicsStories
		global.StoriesUnstartedPercent = global.StoriesUnstarted * 100 / totalEpicsStories
		global.StoriesStartedPercent = global.StoriesStarted * 100 / totalEpicsStories
		global.StoriesBacklogPercent = global.StoriesBacklog * 100 / totalEpicsStories
	}

	totalEpicsEstimateStories := global.EstimateUnstarted + global.EstimateStarted + global.EstimateDone + global.EstimateBacklog
	if totalEpicsEstimateStories != 0 {
		global.EstimateDonePercent = global.EstimateDone * 100 / totalEpicsEstimateStories
		global.EstimateUnstartedPercent = global.EstimateUnstarted * 100 / totalEpicsEstimateStories
		global.EstimateStartedPercent = global.EstimateStarted * 100 / totalEpicsEstimateStories
		global.EstimateBacklogPercent = global.EstimateBacklog * 100 / totalEpicsEstimateStories
	}

	return global
}
