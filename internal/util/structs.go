package util

import (
	"cmp"
	"fmt"
	"log"
	"slices"
	"strconv"
)

var FINAL_HEADERS = []string{
	"year",
	"order",
	"projection",
	"player",
	"age",
	"pos",
	"g",
	"gs",
	"rush_att",
	"rush_yds",
	"rush_td",
	"rush_1d",
	"targets",
	"rec",
	"rec_yds",
	"rec_td",
	"rec_1d",
	"touches",
	"fumbles",
	"pass_cmp",
	"pass_att",
	"pass_yds",
	"pass_td",
	"pass_int",
	"pass_1d",
	"times sacked",
	"pass_sacked_yds",
	"pass_long",
	"rush_long",
	"rec_long",
	"pass_cmp%",
	"pass_att%",
	"pass_yds%",
	"pass_td%",
	"pass_int%",
	"pass_1d%",
	"times sacked%",
	"pass_sacked_yds%",
	"rush_att%",
	"rush_yds%",
	"rush_td%",
	"rush_1d%",
	"targets%",
	"rec%",
	"rec_yds%",
	"rec_td%",
	"rec_1d%",
	"touches%",
	"fumbles%",
	"std_pts",
	"half_ppr_pts",
	"ppr_pts",
	"std_ppg",
	"half_ppr_ppg",
	"ppr_ppg",
	"pos_rank",
}

type Table struct {
	Name      string
	Headers   []string
	Rows      [][]string
	FooterRow []string
}

type TableMap struct {
	Name       string
	Headers    []string
	Dicts      []map[string]string
	FooterDict map[string]string
}

func (t Table) ToMap() TableMap {
	var tableMap TableMap
	tableMap.Name = t.Name
	tableMap.Headers = t.Headers
	tableMap.FooterDict = map[string]string{}
	// convert table rows to data dicts
	rowDicts := []map[string]string{}
	for _, row := range t.Rows {
		rowDict := map[string]string{}
		for i, header := range t.Headers {
			rowDict[header] = row[i]
		}
		rowDicts = append(rowDicts, rowDict)
	}
	tableMap.Dicts = rowDicts
	// convert footer row to dict
	for i, header := range t.Headers {
		tableMap.FooterDict[header] = t.FooterRow[i]
	}
	return tableMap
}

func (m TableMap) ToTable() Table {
	var table Table
	table.Name = m.Name
	table.Headers = m.Headers
	// convert data dicts to table rows
	for _, dict := range m.Dicts {
		row := []string{}
		for _, header := range m.Headers {
			value, exists := dict[header]
			if exists && value != "" {
				row = append(row, value)
			} else if header == "projection" || header == "pos_rank" {
				row = append(row, "")
			} else {
				row = append(row, "0")
			}
		}
		table.Rows = append(table.Rows, row)
	}
	// convert footer dict to table row
	for _, header := range m.Headers {
		value, exists := m.FooterDict[header]
		if exists {
			table.FooterRow = append(table.FooterRow, value)
		} else {
			table.FooterRow = append(table.FooterRow, "")
		}
	}
	return table
}

func (t Table) Sort() Table {
	m := t.ToMap()

	// sort by position then std_pts
	slices.SortFunc(m.Dicts, func(i, j map[string]string) int {
		val1, err := strconv.ParseFloat(i["std_pts"], 64)
		val2, err := strconv.ParseFloat(j["std_pts"], 64)
		if err != nil {
			log.Fatal(err)
		}
		return cmp.Or(
			cmp.Compare(i["pos"], j["pos"]),
			cmp.Compare(val1, val2)*-1,
		)
	})

	return m.ToTable()
}

func MergeTables(tables []Table) Table {
	var mergedTable Table
	var mergedTableMap TableMap
	mergedTableMap.Name = "merged"
	mergedTableMap.FooterDict = map[string]string{}

	if len(tables) > 1 {
		for _, tbl := range tables {
			// append new headers
			for _, header := range tbl.Headers {
				if !slices.Contains(mergedTableMap.Headers, header) {
					mergedTableMap.Headers = append(mergedTableMap.Headers, header)
				}
			}
			// check if record already exists by player name
			tblMap := tbl.ToMap()
			if len(tblMap.Dicts) > 0 {
				for _, record := range tblMap.Dicts {
					recordName := record["player"]
					recordFound := false

					for _, mergedRecord := range mergedTableMap.Dicts {
						// if so, append row data to that record
						if mergedRecord["player"] == recordName {
							recordFound = true
							for _, header := range tbl.Headers {
								mergedRecord[header] = record[header]
							}
							continue
						}
					}
					// if not, append row data as new record
					if !recordFound {
						mergedTableMap.Dicts = append(mergedTableMap.Dicts, record)
					}
				}
			}
			// merge footer row
			for _, header := range tbl.Headers {
				mergedTableMap.FooterDict[header] = tblMap.FooterDict[header]
			}
		}
		mergedTable = mergedTableMap.ToTable()
	} else if len(tables) == 1 {
		mergedTable.Name = mergedTableMap.Name
		mergedTable.Headers = tables[0].Headers
		mergedTable.Rows = tables[0].Rows
		mergedTable.FooterRow = tables[0].FooterRow
	}
	return mergedTable
}

func (t Table) PruneColumns() Table {
	tableMap := t.ToMap()
	tableMap.Headers = FINAL_HEADERS
	prunedTable := tableMap.ToTable()
	return prunedTable
}

func (t Table) AddTeamAndYear(team string, year string) Table {
	tableMap := t.ToMap()
	tableMap.Headers = append(tableMap.Headers, "year")
	for _, dict := range tableMap.Dicts {
		dict["year"] = year
	}
	tableMap.FooterDict["year"] = year
	tableMap.FooterDict["player"] = fmt.Sprintf("%s Totals", team)

	return tableMap.ToTable()
}
