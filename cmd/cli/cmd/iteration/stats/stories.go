package stats

import (
	"fmt"
	"log/slog"
	"os"
	"sort"

	"github.com/lederniermetre/shortcut/pkg/shortcut"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gitlab.com/greyxor/slogor"
)

var storiesCmd = &cobra.Command{
	Use:   "stories",
	Short: "Display statistics about stories",
	Run: func(cmd *cobra.Command, args []string) {
		iterationFlag := cmd.Parent().Parent().PersistentFlags().Lookup("iteration")
		if iterationFlag == nil {
			slog.Error("Can not retrieved iteration flag")
			os.Exit(1)
		}

		iterationName := iterationFlag.Value.String()
		slog.Debug("Search iteration", slog.String("name", iterationName))

		iteration := shortcut.RetrieveIteration(iterationName)

		slog.Info("Iteration retrieved", slog.String("name", *iteration.Name))

		stories := shortcut.StoriesByIteration(*iteration.ID)

		postponedStories := map[string]shortcut.StoryPostponed{}
		epicsStats := map[int64]shortcut.EpicsStats{}
		workflowStates := map[int64]shortcut.WorflowInfo{}
		var totalEstimate int64 = 0

		for _, story := range stories {
			if story.EpicID == nil {
				pterm.Warning.Printfln("Story with no epics: %s", *story.Name)
				continue
			}
			epicID := *story.EpicID
			workflowID := *story.WorkflowID
			workflowStateID := *story.WorkflowStateID

			slog.Debug(
				"Compute story",
				slog.String("name", *story.Name),
				slog.Int64("EpicID", epicID),
			)

			if _, ok := workflowStates[workflowStateID]; !ok {
				workflow := shortcut.GetWorkflow(workflowID)

				for _, wfStates := range workflow.States {
					if *wfStates.ID == workflowStateID {
						slog.Debug("Worflow states", slog.String("worfklow", *workflow.Name), slog.String("name", *wfStates.Type))
						workflowStates[workflowStateID] = shortcut.WorflowInfo{Name: *wfStates.Name, Type: *wfStates.Type}
					}
				}
			}

			if _, ok := epicsStats[epicID]; !ok {
				slog.Debug("Epic stats does not exists", slog.Int64("epicID", epicID))

				epic := shortcut.GetEpic(epicID)

				epicsStats[epicID] = shortcut.EpicsStats{
					Name:       *epic.Name,
					WorkflowID: make(map[int64]map[int64]shortcut.WorkflowStats),
				}
			}

			if ws, ok := epicsStats[epicID].WorkflowID[workflowID][workflowStateID]; ok {
				// If the entry exists, update it
				ws.Count++
				epicsStats[epicID].WorkflowID[workflowID][workflowStateID] = ws
			} else {
				// If the entry doesn't exist, initialize it
				epicsStats[epicID].WorkflowID[workflowID] = make(map[int64]shortcut.WorkflowStats)
				ws := shortcut.WorkflowStats{
					Count: 1,
				}
				epicsStats[epicID].WorkflowID[workflowID][workflowStateID] = ws
			}

			epicsStats[epicID] = shortcut.IncreaseEpicsStoriesCounter(workflowStates[workflowStateID], epicsStats[epicID])

			if story.Estimate == nil {
				pterm.Warning.Printfln("Story with no estimate: %s", *story.Name)
			} else {
				totalEstimate = totalEstimate + *story.Estimate

				epicsStats[epicID] = shortcut.IncreaseEpicsEstimateCounter(workflowStates[workflowStateID], epicsStats[epicID], int(*story.Estimate))
			}

			if len(story.PreviousIterationIds) > 0 {
				postponedStories[*story.Name] = shortcut.StoryPostponed{
					Count:  len(story.PreviousIterationIds),
					Url:    *story.AppURL,
					Status: workflowStates[workflowStateID].Name,
				}
			}
		}

		pterm.DefaultHeader.WithFullWidth().Println("Global iteration stats")

		slog.Info("Stories", slog.Int("count", len(stories)))
		slog.Info("Estimate total", slog.Int("count", int(totalEstimate)))

		globalIterationStats := shortcut.GlobalIterationProgress(epicsStats)

		slog.Info("Stories Unstarted", slog.Int("count", int(globalIterationStats.Unstarted)))
		slog.Info("Stories Started", slog.Int("count", int(globalIterationStats.Started)))
		slog.Info("Stories Done", slog.Int("count", int(globalIterationStats.Done)))

		pterm.DefaultHeader.WithFullWidth().Println("Epics (by stories)")
		epicsTableData := pterm.TableData{{"Epic Name", "Unstarted", "Started", "Done"}}

		for _, epicStat := range epicsStats {
			epicStat = shortcut.SummaryEpicStat(epicStat)
			epicsTableData = append(epicsTableData, []string{epicStat.Name, fmt.Sprintf("%d (%d %%)", epicStat.StoriesUnstarted, epicStat.StoriesUnstartedPercent), fmt.Sprintf("%d (%d %%)", epicStat.StoriesStarted, epicStat.StoriesStartedPercent), fmt.Sprintf("%d (%d %%)", epicStat.StoriesDone, epicStat.StoriesDonePercent)})

			for _, wfState := range epicStat.WorkflowID {
				for wfStateID, stateCount := range wfState {
					slog.Debug("steps", slog.String("state", workflowStates[wfStateID].Name), slog.Int("count", stateCount.Count))
				}
			}
		}

		err := pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(epicsTableData).Render()
		if err != nil {
			slog.Error("Rendering epics (by stories) table", slogor.Err(err))
		}

		pterm.DefaultHeader.WithFullWidth().Println("Epics (by estimates)")
		epicsEstimateTableData := pterm.TableData{{"Epic Name", "Unstarted", "Started", "Done"}}

		for _, epicStart := range epicsStats {
			epicStart = shortcut.SummaryEpicStat(epicStart)
			epicsEstimateTableData = append(epicsEstimateTableData, []string{epicStart.Name, fmt.Sprintf("%d (%d %%)", epicStart.EstimateUnstarted, epicStart.EstimateUnstartedPercent), fmt.Sprintf("%d (%d %%)", epicStart.EstimateStarted, epicStart.EstimateStartedPercent), fmt.Sprintf("%d (%d %%)", epicStart.EstimateDone, epicStart.EstimateDonePercent)})

			for _, wfState := range epicStart.WorkflowID {
				for wfStateID, stateCount := range wfState {
					slog.Debug("steps", slog.String("state", workflowStates[wfStateID].Name), slog.Int("count", stateCount.Count))
				}
			}
		}

		err = pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(epicsEstimateTableData).Render()
		if err != nil {
			slog.Error("Rendering epics (by estimates) table", slogor.Err(err))
		}

		pterm.DefaultHeader.WithFullWidth().Println("Postponed stories")

		keys := make([]string, 0, len(postponedStories))
		for key := range postponedStories {
			keys = append(keys, key)
		}

		sort.SliceStable(keys, func(i, j int) bool {
			return postponedStories[keys[i]].Count > postponedStories[keys[j]].Count
		})

		storiesTableData := pterm.TableData{{"Story Name", "Status", "Nb reported"}}
		for _, key := range keys {
			storiesTableData = append(storiesTableData, []string{key, fmt.Sprint(postponedStories[key].Status), fmt.Sprint(postponedStories[key].Count)})
		}

		err = pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(storiesTableData).Render()
		if err != nil {
			slog.Error("Rendering stories table", slogor.Err(err))
		}
	},
}

func newStoriesCommand() *cobra.Command {
	return storiesCmd
}

func init() {
}
