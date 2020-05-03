package export

import (
	"encoding/csv"
	"os"
	"path/filepath"

	"github.com/SebastianJ/harmony-stats/config"
)

// ExportCSV - exports test suite results as csv
func ExportCSV(fileName string, rows [][]string) (string, error) {
	filePath, err := writeCSVToFile(fileName, rows)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func writeCSVToFile(fileName string, rows [][]string) (string, error) {
	filePath := filepath.Join(config.Configuration.Export.Path, fileName)
	dirPath, _ := filepath.Split(filePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(rows) // calls Flush internally

	if err := csvWriter.Error(); err != nil {
		return "", err
	}

	return filePath, nil
}
