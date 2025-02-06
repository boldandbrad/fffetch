package util

import (
	"encoding/csv"
	"log"
	"os"
)

var OUT_DIRS = []string{"output", "output/fetched_pages", "output/parsed_tables"}

func CreateOutDirs() {
	// create output directories if they don't exist
	for _, dir := range OUT_DIRS {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func WriteCSVFile(filePath string, table Table) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}

	lines := append([][]string{table.Headers}, table.Rows...)
	lines = append(lines, table.FooterRow)
	writer := csv.NewWriter(file)
	err = writer.WriteAll(lines)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteFile(filePath string, contents string) {
	if err := os.WriteFile(filePath, []byte(contents), 0644); err != nil {
		log.Fatal(err)
	}
}
