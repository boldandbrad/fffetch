package pfr

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/boldandbrad/fffetch/internal/util"
)

var PFR_URL = "https://www.pro-football-reference.com/teams"

func FetchPage(teamKey string, year int) string {
	url := fmt.Sprintf("%s/%s/%d.htm", PFR_URL, teamKey, year)
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
	util.WriteFile(filePath, strings.Join(lines, "\n"))
}

func ParsePage(filePath string) []util.Table {
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

	var tables []util.Table
	for _, tableid := range PFR_TABLE_IDS {
		table := parseTable(doc, tableid)
		tables = append(tables, table)
	}
	return tables
}

func parseTable(doc *goquery.Document, tableid string) util.Table {
	var table util.Table
	table.Name = tableid

	doc.Find(fmt.Sprintf("#%s", tableid)).Each(func(i int, tsel *goquery.Selection) {
		if i == 0 {
			// loop through headers
			tsel.Find("th").Each(func(_ int, hsel *goquery.Selection) {
				if hsel != nil {
					key, exists := hsel.Attr("data-stat")
					if exists && key != "ranker" && !strings.Contains(key, "header") {
						table.Headers = append(table.Headers, key)
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
				table.Rows = append(table.Rows, row)
			})

			// grab footer row
			var footerRow []string
			tsel.Find("tfoot").Find("tr").Find("td").Each(func(_ int, csel *goquery.Selection) {
				if csel != nil {
					footerRow = append(footerRow, csel.Text())
				}
			})

			table.FooterRow = footerRow
		}
	})

	return table
}
