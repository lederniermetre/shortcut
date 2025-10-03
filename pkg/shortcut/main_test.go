package shortcut

import (
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestSummaryEpicStat(t *testing.T) {
	testCases := []struct {
		name     string
		actual   EpicsStats
		expected EpicsStats
	}{
		{
			name: "Stardard case",
			actual: EpicsStats{
				Name:              "Epic1",
				StoriesBacklog:    1,
				StoriesUnstarted:  2,
				StoriesStarted:    3,
				StoriesDone:       5,
				EstimateBacklog:   2,
				EstimateUnstarted: 2,
				EstimateStarted:   3,
				EstimateDone:      5,
			},
			expected: EpicsStats{
				Name:                    "Epic1",
				StoriesBacklog:          1,
				StoriesBacklogPercent:   9,
				StoriesUnstarted:        2,
				StoriesUnstartedPercent: 18,
				StoriesStarted:          3,
				StoriesStartedPercent:   27,
				StoriesDone:             5,
				StoriesDonePercent:      45,

				EstimateBacklog:          2,
				EstimateBacklogPercent:   16,
				EstimateUnstarted:        2,
				EstimateUnstartedPercent: 16,
				EstimateStarted:          3,
				EstimateStartedPercent:   25,
				EstimateDone:             5,
				EstimateDonePercent:      41,
			},
		},
		{
			name: "Empty Epic",
			actual: EpicsStats{
				Name:                    "Epic2",
				StoriesBacklog:          0,
				StoriesBacklogPercent:   0,
				StoriesUnstarted:        0,
				StoriesUnstartedPercent: 0,
				StoriesStarted:          0,
				StoriesStartedPercent:   0,
				StoriesDone:             0,
				StoriesDonePercent:      0,

				EstimateBacklog:          0,
				EstimateBacklogPercent:   0,
				EstimateUnstarted:        0,
				EstimateUnstartedPercent: 0,
				EstimateStarted:          0,
				EstimateStartedPercent:   0,
				EstimateDone:             0,
				EstimateDonePercent:      0,
			},
			expected: EpicsStats{
				Name:                    "Epic2",
				StoriesBacklog:          0,
				StoriesBacklogPercent:   0,
				StoriesUnstarted:        0,
				StoriesUnstartedPercent: 0,
				StoriesStarted:          0,
				StoriesStartedPercent:   0,
				StoriesDone:             0,
				StoriesDonePercent:      0,

				EstimateBacklog:          0,
				EstimateBacklogPercent:   0,
				EstimateUnstarted:        0,
				EstimateUnstartedPercent: 0,
				EstimateStarted:          0,
				EstimateStartedPercent:   0,
				EstimateDone:             0,
				EstimateDonePercent:      0,
			},
		},
	}

	for _, tc := range testCases {
		result := SummaryEpicStat(tc.actual)
		assert.EqualValuesf(t, tc.expected, result, "Case '%s' failed", tc.name)
	}
}

func TestIncreaseEpicsCounterMultiType(t *testing.T) {
	// Create a sample EpicsStats instance
	epicStats := EpicsStats{
		Name:             "Sample Epic",
		StoriesUnstarted: 10,
		StoriesStarted:   20,
		StoriesDone:      30,
		WorkflowID:       map[int64]map[int64]WorkflowStats{},
	}

	// Define test cases with different types
	testCases := []struct {
		workflowInfo WorkflowInfo
		expected    int // Expected value for the corresponding counter after calling IncreaseEpicsCounter
	}{
		{WorkflowInfo{Name: "Workflow1", Type: "started"}, epicStats.StoriesStarted + 1},
		{WorkflowInfo{Name: "Workflow2", Type: "unstarted"}, epicStats.StoriesUnstarted + 1},
		{WorkflowInfo{Name: "Workflow3", Type: "done"}, epicStats.StoriesDone + 1},
		{WorkflowInfo{Name: "Workflow4", Type: "backlog"}, epicStats.StoriesBacklog + 1},
		{WorkflowInfo{Name: "UnknownWorkflow", Type: "unknown"}, epicStats.StoriesDone}, // No change for unknown type
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Call the function
		result := IncreaseEpicsStoriesCounter(tc.workflowInfo, epicStats)

		// Check if the corresponding counter is increased
		switch tc.workflowInfo.Type {
		case "started":
			if result.StoriesStarted != tc.expected {
				t.Errorf("For %s, Expected Started counter to be %d, but got %d", tc.workflowInfo.Type, tc.expected, result.StoriesStarted)
			}
		case "unstarted":
			if result.StoriesUnstarted != tc.expected {
				t.Errorf("For %s, Expected Unstarted counter to be %d, but got %d", tc.workflowInfo.Type, tc.expected, result.StoriesUnstarted)
			}
		case "done":
			if result.StoriesDone != tc.expected {
				t.Errorf("For %s, Expected Done counter to be %d, but got %d", tc.workflowInfo.Type, tc.expected, result.StoriesDone)
			}
		case "backlog":
			if result.StoriesBacklog != tc.expected {
				t.Errorf("For %s, Expected Backlog counter to be %d, but got %d", tc.workflowInfo.Type, tc.expected, result.StoriesBacklog)
			}
		default:
			// No change expected for unknown type
			if result.StoriesDone != epicStats.StoriesDone || result.StoriesUnstarted != epicStats.StoriesUnstarted || result.StoriesStarted != epicStats.StoriesStarted {
				t.Errorf("For %s, Expected no change, but got %+v", tc.workflowInfo.Type, result)
			}
		}
	}
}

func TestIncreaseEpicsEstimateCounterMultiType(t *testing.T) {
	// Create a sample EpicsStats instance
	epicStats := EpicsStats{
		Name:              "Sample Epic",
		EstimateUnstarted: 10,
		EstimateStarted:   20,
		EstimateDone:      30,
		WorkflowID:        map[int64]map[int64]WorkflowStats{},
	}

	// Define test cases with different types
	testCases := []struct {
		workflowInfo WorkflowInfo
		expected     int
		estimate     int
	}{
		{WorkflowInfo{Name: "Workflow1", Type: "started"}, epicStats.EstimateStarted + 10, 10},
		{WorkflowInfo{Name: "Workflow2", Type: "unstarted"}, epicStats.EstimateUnstarted + 11, 11},
		{WorkflowInfo{Name: "Workflow3", Type: "done"}, epicStats.EstimateDone + 31, 31},
		{WorkflowInfo{Name: "Workflow4", Type: "backlog"}, epicStats.EstimateBacklog + 5, 5},
		{WorkflowInfo{Name: "UnknownWorkflow", Type: "unknown"}, epicStats.EstimateDone, 30},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Call the function
		result := IncreaseEpicsEstimateCounter(tc.workflowInfo, epicStats, tc.estimate)

		// Check if the corresponding counter is increased
		switch tc.workflowInfo.Type {
		case "started":
			if result.EstimateStarted != tc.expected {
				t.Errorf("For %s, Expected Started counter to be %d, but got %d", tc.workflowInfo.Type, tc.expected, result.EstimateStarted)
			}
		case "unstarted":
			if result.EstimateUnstarted != tc.expected {
				t.Errorf("For %s, Expected Unstarted counter to be %d, but got %d", tc.workflowInfo.Type, tc.expected, result.EstimateUnstarted)
			}
		case "done":
			if result.EstimateDone != tc.expected {
				t.Errorf("For %s, Expected Done counter to be %d, but got %d", tc.workflowInfo.Type, tc.expected, result.EstimateDone)
			}
		case "backlog":
			if result.EstimateBacklog != tc.expected {
				t.Errorf("For %s, Expected Backlog counter to be %d, but got %d", tc.workflowInfo.Type, tc.expected, result.EstimateBacklog)
			}
		default:
			// No change expected for unknown type
			if result.EstimateDone != epicStats.EstimateDone || result.EstimateUnstarted != epicStats.EstimateUnstarted || result.EstimateStarted != epicStats.EstimateStarted {
				t.Errorf("For %s, Expected no change, but got %+v", tc.workflowInfo.Type, result)
			}
		}
	}
}

func TestOrdererOwnersUUID(t *testing.T) {
	ownersUUID := map[strfmt.UUID]int64{
		strfmt.UUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"): 5,
		strfmt.UUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"): 3,
		strfmt.UUID("cccccccc-cccc-cccc-cccc-cccccccccccc"): 8,
	}

	expectedResult := []OwnerStats{
		{strfmt.UUID("cccccccc-cccc-cccc-cccc-cccccccccccc"), 8},
		{strfmt.UUID("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), 5},
		{strfmt.UUID("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"), 3},
	}

	result := OrdererOwnersUUID(ownersUUID)

	// Check if the result matches the expected result
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected: %v, Got: %v", expectedResult, result)
	}
}

func TestComputeEpicGlobalStat(t *testing.T) {
	global := GlobalEpicStats{}

	epicEmpty := EpicsStats{
		StoriesBacklog:   0,
		StoriesUnstarted: 0,
		StoriesStarted:   0,
		StoriesDone:      0,

		EstimateBacklog:   0,
		EstimateUnstarted: 0,
		EstimateStarted:   0,
		EstimateDone:      0,
	}

	expectedResultEmpty := GlobalEpicStats{
		StoriesBacklog:          0,
		StoriesBacklogPercent:   0,
		StoriesUnstarted:        0,
		StoriesUnstartedPercent: 0,
		StoriesStarted:          0,
		StoriesStartedPercent:   0,
		StoriesDone:             0,
		StoriesDonePercent:      0,

		EstimateBacklog:          0,
		EstimateBacklogPercent:   0,
		EstimateUnstarted:        0,
		EstimateStarted:          0,
		EstimateDone:             0,
		EstimateUnstartedPercent: 0,
		EstimateStartedPercent:   0,
		EstimateDonePercent:      0,
	}

	global = ComputeEpicGlobalStat(global, epicEmpty)

	if global != expectedResultEmpty {
		t.Errorf("Expected %+v, but got %+v", expectedResultEmpty, global)
	}

	epics := []EpicsStats{
		{
			StoriesBacklog:   1,
			StoriesUnstarted: 2,
			StoriesStarted:   3,
			StoriesDone:      4,

			EstimateUnstarted: 5,
			EstimateStarted:   6,
			EstimateDone:      7,
			EstimateBacklog:   1,
		},
		{
			StoriesBacklog:   1,
			StoriesUnstarted: 4,
			StoriesStarted:   1,
			StoriesDone:      2,

			EstimateBacklog:   1,
			EstimateUnstarted: 1,
			EstimateStarted:   2,
			EstimateDone:      8,
		},
	}

	expectedResult := GlobalEpicStats{
		StoriesBacklog:          2,
		StoriesBacklogPercent:   11,
		StoriesUnstarted:        6,
		StoriesUnstartedPercent: 33,
		StoriesStarted:          4,
		StoriesStartedPercent:   22,
		StoriesDone:             6,
		StoriesDonePercent:      33,

		EstimateBacklog:          2,
		EstimateBacklogPercent:   6,
		EstimateUnstarted:        6,
		EstimateUnstartedPercent: 19,
		EstimateStarted:          8,
		EstimateStartedPercent:   25,
		EstimateDone:             15,
		EstimateDonePercent:      48,
	}

	for _, epic := range epics {
		global = ComputeEpicGlobalStat(global, epic)
	}

	assert.EqualValuesf(t, expectedResult, global, "%v failed")
}
