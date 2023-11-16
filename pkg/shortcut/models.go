package shortcut

import "github.com/go-openapi/strfmt"

type StoryPostponed struct {
	Count  int
	Url    string
	Status string
}

type EpicsStats struct {
	Name             string
	Unstarted        int
	Started          int
	Done             int
	UnstartedPercent int
	StartedPercent   int
	DonePercent      int
	WorkflowID       map[int64]map[int64]WorkflowStats
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
