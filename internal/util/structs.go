package util

import (
	"slices"
)

type Table struct {
	Name    string
	Headers []string
	Rows    [][]string
}

type TableMap struct {
	Name    string
	Headers []string
	Dicts   []map[string]string
}

func (t Table) ToMap() TableMap {
	var tableMap TableMap
	tableMap.Name = t.Name
	tableMap.Headers = t.Headers
	rowDicts := []map[string]string{}
	for _, row := range t.Rows {
		rowDict := map[string]string{}
		for i, header := range t.Headers {
			rowDict[header] = row[i]
		}
		rowDicts = append(rowDicts, rowDict)
	}
	tableMap.Dicts = rowDicts
	return tableMap
}

func (m TableMap) ToTable() Table {
	var table Table
	table.Name = m.Name
	table.Headers = m.Headers
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
	return table
}

func MergeTables(tables []Table) Table {
	var mergedTableMap TableMap
	mergedTableMap.Name = "merged"

	var mergedTable Table

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
					recordName := record["name_display"]
					recordFound := false
					for _, mergedRecord := range mergedTableMap.Dicts {
						// if so, append row data to that record
						if mergedRecord["name_display"] == recordName {
							recordFound = true
							for _, header := range tbl.Headers {
								mergedRecord[header] = record[header]
							}
							break
						}
					}
					// if not, append row data as new record
					if !recordFound {
						mergedTableMap.Dicts = append(mergedTableMap.Dicts, record)
					}
				}
			}
		}
		mergedTable = mergedTableMap.ToTable()
	} else if len(tables) == 1 {
		mergedTable.Name = mergedTableMap.Name
		mergedTable.Headers = tables[0].Headers
		mergedTable.Rows = tables[0].Rows
	}

	return mergedTable
}
