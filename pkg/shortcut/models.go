package shortcut

import "github.com/go-openapi/strfmt"

type StoryPostponed struct {
	Url    string
	Status string
	Count  int
}

type EpicsStats struct {
	Name string

	StoriesUnstarted        int
	StoriesUnstartedPercent int
	StoriesStarted          int
	StoriesStartedPercent   int
	StoriesDone             int
	StoriesDonePercent      int
	StoriesBacklog          int
	StoriesBacklogPercent   int

	EstimateUnstarted        int
	EstimateUnstartedPercent int
	EstimateStarted          int
	EstimateStartedPercent   int
	EstimateDone             int
	EstimateDonePercent      int
	EstimateBacklog          int
	EstimateBacklogPercent   int

	WorkflowID map[int64]map[int64]WorkflowStats
}

type WorkflowStats struct {
	Count int
}

type WorflowInfo struct {
	Name string
	Type string
}

type OwnerStats struct {
	UUID strfmt.UUID
	Load int64
}

type GlobalEpicStats struct {
	StoriesUnstarted        int
	StoriesUnstartedPercent int
	StoriesStarted          int
	StoriesStartedPercent   int
	StoriesDone             int
	StoriesDonePercent      int
	StoriesBacklog          int
	StoriesBacklogPercent   int

	EstimateUnstarted        int
	EstimateUnstartedPercent int
	EstimateStarted          int
	EstimateStartedPercent   int
	EstimateDone             int
	EstimateDonePercent      int
	EstimateBacklog          int
	EstimateBacklogPercent   int
}
