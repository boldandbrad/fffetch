package util

type Table struct {
	Name    string
	Headers []string
	Rows    [][]string
}

type TableMap struct {
	Name    string
	headers []string
	Dicts   []map[string]string
}

func (t Table) ToMap() TableMap {
	var tableMap TableMap
	tableMap.Name = t.Name
	tableMap.headers = t.Headers
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
	table.Headers = m.headers
	for _, dict := range m.Dicts {
		row := []string{}
		for _, header := range m.headers {
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
