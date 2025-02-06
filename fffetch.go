package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/boldandbrad/fffetch/internal/pfr"
	"github.com/boldandbrad/fffetch/internal/util"
)

var TEAMS = []string{"DET"}
var YEARS = []int{2024}

func main() {

	forceFetch := true

	// create output directories if they don't exist
	util.CreateOutDirs()

	fmt.Println("FFFetching data... 🏈")

	teamsToFetch := map[string]string{}
	if len(TEAMS) == 0 {
		// if no teams provided, fetch data for every team
		teamsToFetch = pfr.PFR_TEAM_KEYS
	} else {
		// fetch data for provided teams only
		for _, team := range TEAMS {
			teamsToFetch[team] = pfr.PFR_TEAM_KEYS[team]
		}
	}

	yearsToFetch := []int{}
	if len(YEARS) == 0 {
		// if no years provided, fetch data for last year
		yearsToFetch = append(yearsToFetch, time.Now().Year()-1)
	} else {
		// fetch data for provided years
		yearsToFetch = YEARS
	}

	// fetch data for each team in each year
	yearCount := 0
	teamCount := 0
	for _, year := range yearsToFetch {
		yearCount += 1
		teamCount = 0
		for team, team_key := range teamsToFetch {
			teamCount += 1

			// fetch and despoof page
			fetchFilePath := fmt.Sprintf("output/fetched_pages/%s_%d.html", team, year)
			_, err := os.Stat(fetchFilePath)
			if errors.Is(err, os.ErrNotExist) || forceFetch == true {
				pageString := pfr.FetchPage(team_key, year)
				util.WriteFile(fetchFilePath, pageString)
				pfr.DespoofPage(fetchFilePath)
				fmt.Printf("    > Fetched %s %d ⬇️\n", team, year)

				// sleep to avoid rate limiting, unless we've already fetched the last team in the last year
				if yearCount < len(yearsToFetch) || teamCount < len(teamsToFetch) {
					time.Sleep(time.Millisecond * time.Duration(rand.IntN(2500-2000)+2000))
				}
			} else {
				fmt.Printf("    > Skipped fetching %s %d, already exists.\n", team, year)
			}

			// parse page table data
			tables := pfr.ParsePage(fetchFilePath)
			for _, table := range tables {
				csvFilePath := fmt.Sprintf("output/parsed_tables/%s_%d_%s.csv", team, year, table.Name)
				util.WriteCSVFile(csvFilePath, table)
			}

			// merge table data
			mergedTable := util.MergeTables(tables)
			csvFilePath := fmt.Sprintf("output/parsed_tables/%s_%d_%s.csv", team, year, mergedTable.Name)
			util.WriteCSVFile(csvFilePath, mergedTable)

			// perform calculations

			// write output to file
		}
	}

	fmt.Println("FFFetching complete! 🌟")
}
