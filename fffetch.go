package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/boldandbrad/fffetch/internal"
)

var TEAMS = []string{"DET"}
var YEARS = []int{2024}

type Table struct {
	name    string
	headers []string
	rows    [][]string
}

func FetchPage(teamKey string, year int) string {
	url := fmt.Sprintf("https://www.pro-football-reference.com/teams/%s/%d.htm", teamKey, year)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	resStatus := res.StatusCode
	if resStatus != 200 {
		if resStatus == 429 {
			log.Fatal("Pro Football Reference Rate limit hit, please try again later.")
		} else if resStatus == 404 {
			log.Fatalf("Page not found for %s %d", teamKey, year)
		} else {
			log.Fatalf("Unknown status code %d for %s %d", resStatus, teamKey, year)
		}
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyBytes)
}

func DespoofPage(filePath string) {
	// read file into memory
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// remove start and end comments
	startPattern, _ := regexp.Compile("^[ ]*<!--$")
	endPattern, _ := regexp.Compile("^-->$")
	altEndPattern, _ := regexp.Compile("See more advanced stats here.</a></div>-->$")
	for i, line := range lines {
		if startPattern.MatchString(line) || endPattern.MatchString(line) {
			lines[i] = ""
		} else if altEndPattern.MatchString(line) {
			lines[i] = strings.Replace(line, "-->", "", 1)
		}
	}

	// write file back
	WriteFile(filePath, strings.Join(lines, "\n"))
}

func ParseTable(doc *goquery.Document, tableid string) Table {
	var table Table
	table.name = tableid

	doc.Find(fmt.Sprintf("#%s", tableid)).Each(func(i int, tsel *goquery.Selection) {
		if i == 0 {
			// loop through headers
			tsel.Find("th").Each(func(_ int, hsel *goquery.Selection) {
				if hsel != nil {
					key, exists := hsel.Attr("data-stat")
					if exists && key != "ranker" && !strings.Contains(key, "header") {
						table.headers = append(table.headers, key)
					}
				}
			})

			// loop through rows
			tsel.Find("tbody").Find("tr").Each(func(index int, rsel *goquery.Selection) {
				var row []string

				// loop through cells
				rsel.Find("td").Each(func(_ int, csel *goquery.Selection) {
					if csel != nil {
						row = append(row, csel.Text())
					}
				})
				table.rows = append(table.rows, row)
			})
		}
	})

	fmt.Println(table.rows)
	return table
}

func ParsePage(filePath string) []Table {
	// read file into memory
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	var tables []Table
	for _, tableid := range internal.PFR_TABLE_IDS {
		table := ParseTable(doc, tableid)
		tables = append(tables, table)
	}
	return tables
}

func WriteFile(filePath string, contents string) {
	if err := os.WriteFile(filePath, []byte(contents), 0644); err != nil {
		log.Fatal(err)
	}
}

func WriteCSVFile(filePath string, table Table) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	lines := append([][]string{table.headers}, table.rows...)
	writer := csv.NewWriter(file)
	err = writer.WriteAll(lines)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	forceFetch := true

	// create output directories if they don't exist
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		if err := os.Mkdir("output", 0755); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat("output/fetched_pages"); os.IsNotExist(err) {
		if err := os.Mkdir("output/fetched_pages", 0755); err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat("output/parsed_tables"); os.IsNotExist(err) {
		if err := os.Mkdir("output/parsed_tables", 0755); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("FFFetching data... 🏈")

	teamsToFetch := map[string]string{}
	if len(TEAMS) == 0 {
		// if no teams provided, fetch data for every team
		teamsToFetch = internal.PFR_TEAM_KEYS
	} else {
		// fetch data for provided teams only
		for _, team := range TEAMS {
			teamsToFetch[team] = internal.PFR_TEAM_KEYS[team]
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
				pageString := FetchPage(team_key, year)
				WriteFile(fetchFilePath, pageString)
				DespoofPage(fetchFilePath)
				fmt.Printf("    > Fetched %s %d ⬇️\n", team, year)
			} else {
				fmt.Printf("    > Skipped fetching %s %d, already exists.\n", team, year)
			}

			// parse page table data
			tables := ParsePage(fetchFilePath)
			for _, table := range tables {
				csvFilePath := fmt.Sprintf("output/parsed_tables/%s_%d_%s.csv", team, year, table.name)
				WriteCSVFile(csvFilePath, table)
			}

			// perform calculations

			// write output to file

			// sleep to avoid rate limiting, unless we've already fetched the last team in the last year
			if yearCount < len(yearsToFetch) || teamCount < len(teamsToFetch) {
				time.Sleep(time.Millisecond * time.Duration(rand.IntN(2500-2000)+2000))
			}
		}
	}

	fmt.Println("FFFetching complete! 🌟")
}
