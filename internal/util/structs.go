package util

import (
	"slices"
)

var OUT_HEADERS = []string{
	"player",
	"age",
	"pos",
	"g",
	"gs",
	"pass_cmp",
	"pass_att",
	"pass_yds",
	"pass_td",
	"pass_int",
	"pass_sacked_yds",
	"rush_att",
	"rush_yds",
	"rush_td",
	"rush_1d",
	"rush_1d%",
	"targets",
	"rec",
	"rec_yds",
	"rec_td",
	"rec_1d",
	"rec_1d%",
	"touches",
	"fumbles",
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
			if exists {
				row = append(row, value)
			} else {
				row = append(row, "")
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
	tableMap.Headers = OUT_HEADERS
	prunedTable := tableMap.ToTable()
	return prunedTable
}
