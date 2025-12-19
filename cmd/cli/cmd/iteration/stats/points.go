package stats

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lederniermetre/shortcut/pkg/shortcut"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type IterationPoints struct {
	ID     int64
	Name   string
	Points int64
}

var pointsCmd = &cobra.Command{
	Use:   "points",
	Short: "Display iterations with most and least points",
	Run: func(cmd *cobra.Command, args []string) {
		queryFlag := cmd.Parent().PersistentFlags().Lookup("query")
		if queryFlag == nil {
			slog.Error("Can not retrieved query flag")
			os.Exit(1)
		}
		limitFlag, err := cmd.Parent().Flags().GetInt("limit")
		if err != nil {
			slog.Error("Can not retrieved limit flag", slog.Any("error", err))
			os.Exit(1)
		}

		shortcutQuery := queryFlag.Value.String()
		slog.Debug("Search", slog.String("name", shortcutQuery))

		iterations, err := shortcut.RetrieveIterations(shortcutQuery, limitFlag, "")
		if err != nil {
			slog.Error("Cannot retrieve iterations", slog.Any("error", err))
			os.Exit(1)
		}

		if len(iterations) == 0 {
			pterm.Warning.Println("No iterations found")
			return
		}

		var iterationPointsList []IterationPoints
		var maxPoints IterationPoints
		var minPoints IterationPoints
		minPoints.Points = -1 // Initialize to -1 to detect first iteration

		for _, iteration := range iterations {
			slog.Info("Iteration retrieved", slog.String("name", *iteration.Name))

			stories, err := shortcut.StoriesByIteration(*iteration.ID)
			if err != nil {
				slog.Error("Cannot retrieve stories", slog.Any("error", err))
				os.Exit(1)
			}

			var totalPoints int64 = 0
			for _, story := range stories {
				if story.Archived != nil && *story.Archived {
					slog.Debug("Story archived, skipping", slog.String("name", *story.Name))
					continue
				}

				if story.Estimate != nil {
					totalPoints += *story.Estimate
				}
			}

			iterPoint := IterationPoints{
				ID:     *iteration.ID,
				Name:   *iteration.Name,
				Points: totalPoints,
			}
			iterationPointsList = append(iterationPointsList, iterPoint)

			// Track max and min
			if totalPoints > maxPoints.Points {
				maxPoints = iterPoint
			}

			if minPoints.Points == -1 || totalPoints < minPoints.Points {
				minPoints = iterPoint
			}
		}

		// Display all iterations
		pterm.DefaultHeader.WithFullWidth().Println("Iteration Points")
		allIterationsTable := pterm.TableData{{"Iteration Name", "Total Points"}}
		for _, iter := range iterationPointsList {
			allIterationsTable = append(allIterationsTable, []string{iter.Name, fmt.Sprintf("%d", iter.Points)})
		}

		err = pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(allIterationsTable).Render()
		if err != nil {
			slog.Error("Rendering iterations table", slog.Any("error", err))
		}

		// Display max and min
		pterm.DefaultHeader.WithFullWidth().Println("Iteration Statistics")
		statsTable := pterm.TableData{
			{"Statistic", "Iteration Name", "Total Points"},
			{pterm.FgGreen.Sprint("Most Points"), maxPoints.Name, fmt.Sprintf("%d", maxPoints.Points)},
			{pterm.FgRed.Sprint("Least Points"), minPoints.Name, fmt.Sprintf("%d", minPoints.Points)},
		}

		err = pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(statsTable).Render()
		if err != nil {
			slog.Error("Rendering statistics table", slog.Any("error", err))
		}
	},
}

func newPointsCommand() *cobra.Command {
	return pointsCmd
}

func init() {
}
