package cmd

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/boldandbrad/fffetch/internal/calc"
	"github.com/boldandbrad/fffetch/internal/pfr"
	"github.com/boldandbrad/fffetch/internal/util"
	"github.com/boldandbrad/fffetch/pkg/tea"
	"github.com/spf13/cobra"
)

var (
	teams []string
	years []int
	force bool
	show  bool
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch fantasy football data",
	Long:  "Fetch and process fantasy football data from Pro Football Reference",
	Run: func(cmd *cobra.Command, args []string) {
		runFetch()
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.Flags().StringSliceVarP(&teams, "team", "t", []string{}, "Teams to fetch (e.g., KC, BUF, PHI). Defaults to all teams")
	fetchCmd.Flags().IntSliceVarP(&years, "year", "y", []int{}, "Years to fetch (e.g., 2023, 2024). Defaults to previous year")
	fetchCmd.Flags().BoolVarP(&force, "force", "f", false, "Force re-fetch existing data")
}

func runFetch() {
	forceFetch := force

	util.CreateOutDirs()

	teamsToFetch := map[string]string{}
	if len(teams) == 0 {
		teamsToFetch = pfr.PFR_TEAM_KEYS
	} else {
		for _, team := range teams {
			if key, exists := pfr.PFR_TEAM_KEYS[team]; exists {
				teamsToFetch[team] = key
			} else {
				fmt.Printf("Invalid team: %s\n", team)
			}
		}
		if len(teamsToFetch) == 0 {
			fmt.Println("No valid teams provided")
			os.Exit(1)
		}
	}

	yearsToFetch := []int{}
	if len(years) == 0 {
		yearsToFetch = append(yearsToFetch, time.Now().Year()-1)
	} else {
		yearsToFetch = years
	}

	totalTasks := len(yearsToFetch) * len(teamsToFetch)
	if totalTasks == 0 {
		fmt.Println("No teams or years to fetch")
		return
	}

	p := tea.NewProgram(totalTasks)
	p.Start()

	done := make(chan struct{})

	go func() {
		yearCount := 0
		teamCount := 0
		for _, year := range yearsToFetch {
			yearCount += 1
			teamCount = 0
			for team, team_key := range teamsToFetch {
				teamCount += 1

				fetchFilePath := fmt.Sprintf("output/fetched_pages/%s_%d.html", team, year)
				_, err := os.Stat(fetchFilePath)
				if errors.Is(err, os.ErrNotExist) || forceFetch {
					pageString := pfr.FetchPage(team_key, year)
					util.WriteFile(fetchFilePath, pageString)

					tables := pfr.ParsePage(fetchFilePath)
					for _, table := range tables {
						csvFilePath := fmt.Sprintf("output/parsed_tables/%s_%d_%s.csv", team, year, table.Name)
						util.WriteCSVFile(csvFilePath, table)
					}

					mergedTable := util.MergeTables(tables)
					csvFilePath := fmt.Sprintf("output/parsed_tables/%s_%d_%s.csv", team, year, mergedTable.Name)
					util.WriteCSVFile(csvFilePath, mergedTable)

					statTable := calc.CalcAdvStats(mergedTable)
					statTable = calc.CalcFFStats(statTable)

					updatedTable := statTable.AddTeamAndYear(team, strconv.Itoa(year))
					updatedTable = updatedTable.Sort()
					prunedTable := updatedTable.PruneColumns()
					csvFilePath = fmt.Sprintf("output/final/%s_%d.csv", team, year)
					util.WriteCSVFile(csvFilePath, prunedTable)

					p.Update(tea.TaskResult{Team: team, Year: year, Success: true})

					if yearCount < len(yearsToFetch) || teamCount < len(teamsToFetch) {
						time.Sleep(time.Millisecond * time.Duration(rand.IntN(2500-2000)+2000))
					}
				} else {
					p.Update(tea.TaskResult{Team: team, Year: year, Success: false})
				}
			}
		}
		close(done)
	}()

	<-done
	p.Quit()
	time.Sleep(100 * time.Millisecond)
}
