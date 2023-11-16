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

		allStories := shortcut.StoriesByIteration(*iteration.ID)

		postponedStories := map[string]shortcut.StoryPostponed{}
		epicsStats := map[int64]shortcut.EpicsStats{}
		workflowStates := map[int64]shortcut.WorflowInfo{}
		var totalEstimate int64 = 0

		for _, story := range allStories {
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

				epicsStats[epicID].WorkflowID[workflowID] = make(map[int64]shortcut.WorkflowStats)
				ws := shortcut.WorkflowStats{
					Count: 1,
				}
				epicsStats[epicID].WorkflowID[workflowID][workflowStateID] = ws
			} else {
				if ws, ok := epicsStats[epicID].WorkflowID[workflowID][workflowStateID]; ok {
					// If the entry exists, update it
					ws.Count++
					epicsStats[epicID].WorkflowID[workflowID][workflowStateID] = ws
				} else {
					// If the entry doesn't exist, initialize it
					ws := shortcut.WorkflowStats{
						Count: 1,
					}
					epicsStats[epicID].WorkflowID[workflowID][workflowStateID] = ws
				}
			}

			epicsStats[epicID] = shortcut.IncreaseEpicsCounter(workflowStates[workflowStateID], epicsStats[epicID])

			if story.Estimate == nil {
				pterm.Warning.Printfln("Story with no estimate: %s", *story.Name)
			} else {
				totalEstimate = totalEstimate + *story.Estimate
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

		slog.Info("Number of stories", slog.Int("count", len(allStories)))
		slog.Info("Estimate total", slog.Int("count", int(totalEstimate)))

		pterm.DefaultHeader.WithFullWidth().Println("Epics")
		epicsTableData := pterm.TableData{{"Epic Name", "Unstarted", "Started", "Done"}}

		for _, v := range epicsStats {
			v = shortcut.SummaryEpicStat(v)
			epicsTableData = append(epicsTableData, []string{v.Name, fmt.Sprintf("%d (%d %%)", v.Unstarted, v.UnstartedPercent), fmt.Sprintf("%d (%d %%)", v.Started, v.StartedPercent), fmt.Sprintf("%d (%d %%)", v.Done, v.DonePercent)})

			for _, wfState := range v.WorkflowID {
				for wfStateID, stateCount := range wfState {
					slog.Debug("steps", slog.String("state", workflowStates[wfStateID].Name), slog.Int("count", stateCount.Count))
				}
			}
		}

		err := pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(epicsTableData).Render()
		if err != nil {
			slog.Error("Rendering epics table", slogor.Err(err))
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
		for _, k := range keys {
			storiesTableData = append(storiesTableData, []string{k, fmt.Sprint(postponedStories[k].Status), fmt.Sprint(postponedStories[k].Count)})
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
