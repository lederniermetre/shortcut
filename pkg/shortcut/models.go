package shortcut

type StoryPostponed struct {
	Count  int
	Url    string
	Status string
}

type EpicsStats struct {
	Name       string
	WorkflowID map[int64]map[int64]WorkflowStats
}

type WorkflowStats struct {
	Count int
}
