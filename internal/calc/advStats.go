package calc

import (
	"fmt"
	"log"
	"slices"
	"strconv"

	"github.com/boldandbrad/fffetch/internal/util"
)

var fieldsToPercent = []string{
	"pass_cmp",
	"pass_att",
	"pass_yds",
	"pass_td",
	"pass_int",
	"pass_1d",
	"times sacked",
	"pass_sacked_yds",
	"rush_att",
	"rush_yds",
	"rush_td",
	"targets",
	"rec",
	"rec_yds",
	"rec_td",
	"touches",
	"fumbles",
}

func CalcAdvStats(table util.Table) util.Table {
	tableMap := table.ToMap()

	// add advanced stat headers
	for _, header := range fieldsToPercent {
		adjFieldName := fmt.Sprintf("%s%%", header)
		if !slices.Contains(tableMap.Headers, adjFieldName) {
			tableMap.Headers = append(tableMap.Headers, adjFieldName)
		}
	}

	// calculate advanced stats for each player
	for _, dict := range tableMap.Dicts {
		for _, field := range fieldsToPercent {
			adjFieldName := fmt.Sprintf("%s%%", field)
			teamTotal, err := strconv.ParseFloat(tableMap.FooterDict[field], 32)
			playerVal, err := strconv.ParseFloat(dict[field], 32)
			if err != nil {
				log.Fatal(err)
			}
			percentage := (playerVal / teamTotal)
			dict[adjFieldName] = fmt.Sprintf("%.2f%%", percentage*100)
		}
	}

	return tableMap.ToTable()
}
