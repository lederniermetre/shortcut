package shortcut

import (
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"
)

func TestSummaryEpicStat(t *testing.T) {
	epic := EpicsStats{
		Name:                     "Epic1",
		StoriesUnstarted:         2,
		StoriesStarted:           3,
		StoriesDone:              5,
		EstimateUnstarted:        2,
		EstimateStarted:          3,
		EstimateDone:             5,
		StoriesUnstartedPercent:  0,
		StoriesStartedPercent:    0,
		StoriesDonePercent:       0,
		EstimateUnstartedPercent: 0,
		EstimateStartedPercent:   0,
		EstimateDonePercent:      0,
	}

	result := SummaryEpicStat(epic)

	if result.StoriesUnstartedPercent != 20 ||
		result.StoriesStartedPercent != 30 ||
		result.StoriesDonePercent != 50 ||
		result.EstimateUnstartedPercent != 20 ||
		result.EstimateStartedPercent != 30 ||
		result.EstimateDonePercent != 50 {
		t.Errorf("Percentage calculation error")
	}

	epicZero := EpicsStats{
		Name:                     "Epic2",
		StoriesUnstarted:         0,
		StoriesStarted:           0,
		StoriesDone:              0,
		EstimateUnstarted:        0,
		EstimateStarted:          0,
		EstimateDone:             0,
		StoriesUnstartedPercent:  0,
		StoriesStartedPercent:    0,
		StoriesDonePercent:       0,
		EstimateUnstartedPercent: 0,
		EstimateStartedPercent:   0,
		EstimateDonePercent:      0,
	}

	resultZero := SummaryEpicStat(epicZero)

	if resultZero.StoriesUnstartedPercent != 0 ||
		resultZero.StoriesStartedPercent != 0 ||
		resultZero.StoriesDonePercent != 0 ||
		resultZero.EstimateUnstartedPercent != 0 ||
		resultZero.EstimateStartedPercent != 0 ||
		resultZero.EstimateDonePercent != 0 {
		t.Errorf("Percentage calculation error for zero values")
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
		worflowInfo WorflowInfo
		expected    int // Expected value for the corresponding counter after calling IncreaseEpicsCounter
	}{
		{WorflowInfo{Name: "Workflow1", Type: "started"}, epicStats.StoriesStarted + 1},
		{WorflowInfo{Name: "Workflow2", Type: "unstarted"}, epicStats.StoriesUnstarted + 1},
		{WorflowInfo{Name: "Workflow3", Type: "done"}, epicStats.StoriesDone + 1},
		{WorflowInfo{Name: "UnknownWorkflow", Type: "unknown"}, epicStats.StoriesDone}, // No change for unknown type
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Call the function
		result := IncreaseEpicsStoriesCounter(tc.worflowInfo, epicStats)

		// Check if the corresponding counter is increased
		switch tc.worflowInfo.Type {
		case "started":
			if result.StoriesStarted != tc.expected {
				t.Errorf("For %s, Expected Started counter to be %d, but got %d", tc.worflowInfo.Type, tc.expected, result.StoriesStarted)
			}
		case "unstarted":
			if result.StoriesUnstarted != tc.expected {
				t.Errorf("For %s, Expected Unstarted counter to be %d, but got %d", tc.worflowInfo.Type, tc.expected, result.StoriesUnstarted)
			}
		case "done":
			if result.StoriesDone != tc.expected {
				t.Errorf("For %s, Expected Done counter to be %d, but got %d", tc.worflowInfo.Type, tc.expected, result.StoriesDone)
			}
		default:
			// No change expected for unknown type
			if result.StoriesDone != epicStats.StoriesDone || result.StoriesUnstarted != epicStats.StoriesUnstarted || result.StoriesStarted != epicStats.StoriesStarted {
				t.Errorf("For %s, Expected no change, but got %+v", tc.worflowInfo.Type, result)
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
		worflowInfo WorflowInfo
		expected    int
		estimate    int
	}{
		{WorflowInfo{Name: "Workflow1", Type: "started"}, epicStats.EstimateStarted + 10, 10},
		{WorflowInfo{Name: "Workflow2", Type: "unstarted"}, epicStats.EstimateUnstarted + 11, 11},
		{WorflowInfo{Name: "Workflow3", Type: "done"}, epicStats.EstimateDone + 31, 31},
		{WorflowInfo{Name: "UnknownWorkflow", Type: "unknown"}, epicStats.EstimateDone, 30},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Call the function
		result := IncreaseEpicsEstimateCounter(tc.worflowInfo, epicStats, tc.estimate)

		// Check if the corresponding counter is increased
		switch tc.worflowInfo.Type {
		case "started":
			if result.EstimateStarted != tc.expected {
				t.Errorf("For %s, Expected Started counter to be %d, but got %d", tc.worflowInfo.Type, tc.expected, result.EstimateStarted)
			}
		case "unstarted":
			if result.EstimateUnstarted != tc.expected {
				t.Errorf("For %s, Expected Unstarted counter to be %d, but got %d", tc.worflowInfo.Type, tc.expected, result.EstimateUnstarted)
			}
		case "done":
			if result.EstimateDone != tc.expected {
				t.Errorf("For %s, Expected Done counter to be %d, but got %d", tc.worflowInfo.Type, tc.expected, result.EstimateDone)
			}
		default:
			// No change expected for unknown type
			if result.EstimateDone != epicStats.EstimateDone || result.EstimateUnstarted != epicStats.EstimateUnstarted || result.EstimateStarted != epicStats.EstimateStarted {
				t.Errorf("For %s, Expected no change, but got %+v", tc.worflowInfo.Type, result)
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
		StoriesUnstarted:  0,
		StoriesStarted:    0,
		StoriesDone:       0,
		EstimateUnstarted: 0,
		EstimateStarted:   0,
		EstimateDone:      0,
	}

	expectedResultEmpty := GlobalEpicStats{
		StoriesUnstarted:         0,
		StoriesStarted:           0,
		StoriesDone:              0,
		EstimateUnstarted:        0,
		EstimateStarted:          0,
		EstimateDone:             0,
		StoriesUnstartedPercent:  0,
		StoriesStartedPercent:    0,
		StoriesDonePercent:       0,
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
			StoriesUnstarted:  2,
			StoriesStarted:    3,
			StoriesDone:       4,
			EstimateUnstarted: 5,
			EstimateStarted:   6,
			EstimateDone:      7,
		},
		{
			StoriesUnstarted:  4,
			StoriesStarted:    1,
			StoriesDone:       2,
			EstimateUnstarted: 1,
			EstimateStarted:   2,
			EstimateDone:      8,
		},
	}

	expectedResult := GlobalEpicStats{
		StoriesUnstarted:         6,
		StoriesStarted:           4,
		StoriesDone:              6,
		EstimateUnstarted:        6,
		EstimateStarted:          8,
		EstimateDone:             15,
		StoriesUnstartedPercent:  37,
		StoriesStartedPercent:    25,
		StoriesDonePercent:       37,
		EstimateUnstartedPercent: 20,
		EstimateStartedPercent:   27,
		EstimateDonePercent:      51,
	}

	for _, epic := range epics {
		global = ComputeEpicGlobalStat(global, epic)
	}

	if global != expectedResult {
		t.Errorf("Expected %+v, but got %+v", expectedResult, global)
	}
}
