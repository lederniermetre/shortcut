package shortcut

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestStoryPostponedStruct(t *testing.T) {
	sp := StoryPostponed{
		Url:    "https://example.com/story/123",
		Status: "in_progress",
		Count:  3,
	}

	assert.Equal(t, "https://example.com/story/123", sp.Url)
	assert.Equal(t, "in_progress", sp.Status)
	assert.Equal(t, 3, sp.Count)
}

func TestEpicsStatsStruct(t *testing.T) {
	es := EpicsStats{
		Name:                     "Epic 1",
		WorkflowID:               make(map[int64]map[int64]WorkflowStats),
		StoriesUnstarted:         5,
		StoriesUnstartedPercent:  50,
		StoriesStarted:           3,
		StoriesStartedPercent:    30,
		StoriesDone:              2,
		StoriesDonePercent:       20,
		EstimateUnstarted:        10,
		EstimateUnstartedPercent: 40,
		EstimateStarted:          8,
		EstimateStartedPercent:   32,
		EstimateDone:             7,
		EstimateDonePercent:      28,
	}

	assert.Equal(t, "Epic 1", es.Name)
	assert.Equal(t, 5, es.StoriesUnstarted)
	assert.Equal(t, 50, es.StoriesUnstartedPercent)
}

func TestWorkflowInfoStruct(t *testing.T) {
	wf := WorkflowInfo{
		Name: "Development",
		Type: "started",
	}

	assert.Equal(t, "Development", wf.Name)
	assert.Equal(t, "started", wf.Type)
}

func TestOwnerStatsStruct(t *testing.T) {
	uuid := strfmt.UUID("test-uuid-123")
	os := OwnerStats{
		UUID: uuid,
		Load: 42,
	}

	assert.Equal(t, uuid, os.UUID)
	assert.Equal(t, int64(42), os.Load)
}

func TestGlobalEpicStatsStruct(t *testing.T) {
	ges := GlobalEpicStats{
		StoriesUnstarted:         10,
		StoriesUnstartedPercent:  40,
		StoriesStarted:           8,
		StoriesStartedPercent:    32,
		StoriesDone:              7,
		StoriesDonePercent:       28,
		EstimateUnstarted:        20,
		EstimateUnstartedPercent: 45,
		EstimateStarted:          15,
		EstimateStartedPercent:   34,
		EstimateDone:             10,
		EstimateDonePercent:      21,
	}

	assert.Equal(t, 10, ges.StoriesUnstarted)
	assert.Equal(t, 40, ges.StoriesUnstartedPercent)
	assert.Equal(t, 20, ges.EstimateUnstarted)
}

func TestWorkflowStatsStruct(t *testing.T) {
	ws := WorkflowStats{
		Count: 5,
	}

	assert.Equal(t, 5, ws.Count)
}
