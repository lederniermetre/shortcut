package shortcut

import "github.com/go-openapi/strfmt"

type StoryPostponed struct {
	Count  int
	Url    string
	Status string
}

type EpicsStats struct {
	Name string

	StoriesUnstarted        int
	StoriesUnstartedPercent int
	StoriesStarted          int
	StoriesStartedPercent   int
	StoriesDone             int
	StoriesDonePercent      int

	EstimateUnstarted        int
	EstimateUnstartedPercent int
	EstimateStarted          int
	EstimateStartedPercent   int
	EstimateDone             int
	EstimateDonePercent      int

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
