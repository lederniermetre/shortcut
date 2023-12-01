package shortcut

import (
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"
)

func TestSummaryEpicStat(t *testing.T) {
	epic := EpicsStats{
		Name:             "Sample Epic",
		StoriesUnstarted: 10,
		StoriesStarted:   20,
		StoriesDone:      30,
		WorkflowID:       map[int64]map[int64]WorkflowStats{},
	}

	// Call the function
	result := SummaryEpicStat(epic)

	// Check if the percentages are calculated correctly
	expectedDonePercent := 30 * 100 / (10 + 20 + 30)
	expectedUnstartedPercent := 10 * 100 / (10 + 20 + 30)
	expectedStartedPercent := 20 * 100 / (10 + 20 + 30)

	if result.StoriesDonePercent != expectedDonePercent {
		t.Errorf("Expected DonePercent to be %d, but got %d", expectedDonePercent, result.StoriesDonePercent)
	}

	if result.StoriesUnstartedPercent != expectedUnstartedPercent {
		t.Errorf("Expected UnstartedPercent to be %d, but got %d", expectedUnstartedPercent, result.StoriesUnstartedPercent)
	}

	if result.StoriesStartedPercent != expectedStartedPercent {
		t.Errorf("Expected StartedPercent to be %d, but got %d", expectedStartedPercent, result.StoriesStartedPercent)
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
		result := IncreaseEpicsCounter(tc.worflowInfo, epicStats)

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
