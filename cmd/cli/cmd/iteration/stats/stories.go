package stats

import (
	"fmt"
	"log/slog"
	"os"
	"sort"

	"github.com/lederniermetre/shortcut/pkg/shortcut"
	"github.com/spf13/cobra"
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
		workflowStates := map[int64]string{}
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
						workflowStates[workflowStateID] = *wfStates.Name
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

			if story.Estimate == nil {
				slog.Warn("Story with no estimate", slog.String("name", *story.Name))
			} else {
				totalEstimate = totalEstimate + *story.Estimate
			}

			if len(story.PreviousIterationIds) > 0 {
				postponedStories[*story.Name] = shortcut.StoryPostponed{
					Count:  len(story.PreviousIterationIds),
					Url:    *story.AppURL,
					Status: workflowStates[workflowStateID],
				}
			}
		}

		slog.Info("===== Global stats =====")

		slog.Info("Number of stories", slog.Int("count", len(allStories)))
		slog.Info("Estimate total", slog.Int("count", int(totalEstimate)))

		slog.Info("===== Epics =====")

		for _, v := range epicsStats {
			for _, wfState := range v.WorkflowID {
				for wfStateID, stateCount := range wfState {
					slog.Info("Epic stats", slog.String("name", v.Name), slog.String("state", workflowStates[wfStateID]), slog.Int("count", stateCount.Count))
				}
			}
		}

		slog.Info("===== Postponed stories =====")

		keys := make([]string, 0, len(postponedStories))
		for key := range postponedStories {
			keys = append(keys, key)
		}

		sort.SliceStable(keys, func(i, j int) bool {
			return postponedStories[keys[i]].Count > postponedStories[keys[j]].Count
		})

		for _, k := range keys {
			slog.Info(fmt.Sprintf("%s has been postponed", k), slog.Int("count", postponedStories[k].Count), slog.String("status", postponedStories[k].Status))
			slog.Debug("access story", slog.String("url", postponedStories[k].Url))
		}
	},
}

func newStoriesCommand() *cobra.Command {
	return storiesCmd
}

func init() {
}
