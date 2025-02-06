package calc

import (
	"fmt"
	"log"
	"slices"
	"strconv"

	"github.com/boldandbrad/fffetch/internal/util"
)

var fantasyFootballFields = []string{
	"std_pts",
	"half_ppr_pts",
	"ppr_pts",
	"std_ppg",
	"half_ppr_ppg",
	"ppr_ppg",
}

func CalcFFStats(table util.Table) util.Table {
	tableMap := table.ToMap()

	// add fantasy football stat headers
	for _, header := range fantasyFootballFields {
		if !slices.Contains(tableMap.Headers, header) {
			tableMap.Headers = append(tableMap.Headers, header)
		}
	}

	// calculate fantasy football stats for each player
	for _, dict := range tableMap.Dicts {
		games, err := strconv.Atoi(dict["g"])
		rushYds, err := strconv.Atoi(dict["rush_yds"])
		rushTD, err := strconv.Atoi(dict["rush_td"])
		rec, err := strconv.Atoi(dict["rec"])
		recYds, err := strconv.Atoi(dict["rec_yds"])
		recTD, err := strconv.Atoi(dict["rec_td"])
		fumbles, err := strconv.Atoi(dict["fumbles"])
		passYds, err := strconv.Atoi(dict["pass_yds"])
		passTD, err := strconv.Atoi(dict["pass_td"])
		passInt, err := strconv.Atoi(dict["pass_int"])
		if err != nil {
			log.Fatal(err)
		}

		// standard
		stdPts := (float64(rushYds) * 0.1) + (float64(rushTD) * 6) + (float64(recYds) * 0.1) + (float64(recTD) * 6) + (float64(fumbles) * -1) + (float64(passYds) * 0.04) + (float64(passTD) * 4) + (float64(passInt) * -2)
		dict["std_pts"] = fmt.Sprintf("%.2f", stdPts)
		stdPpg := stdPts / float64(games)
		dict["std_ppg"] = fmt.Sprintf("%.2f", stdPpg)

		// half point per reception
		halfPprPts := stdPts + (float64(rec) * 0.5)
		dict["half_ppr_pts"] = fmt.Sprintf("%.2f", halfPprPts)
		halfPprPpg := halfPprPts / float64(games)
		dict["half_ppr_ppg"] = fmt.Sprintf("%.2f", halfPprPpg)

		// point per reception
		pprPts := stdPts + (float64(rec) * 1)
		dict["ppr_pts"] = fmt.Sprintf("%.2f", pprPts)
		pprPpg := halfPprPts / float64(games)
		dict["ppr_ppg"] = fmt.Sprintf("%.2f", pprPpg)
	}

	// TODO: calculate order

	return tableMap.ToTable()
}
