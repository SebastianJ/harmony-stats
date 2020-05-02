package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SebastianJ/harmony-stats/config"
	"github.com/SebastianJ/harmony-stats/utils"
)

var (
	timeFormat string = "2006-01-02 15:04:05 UTC"
)

func generateFileName(theTime time.Time, ext string) string {
	return fmt.Sprintf("validators-%s-UTC.%s", utils.FormattedTimeString(theTime), ext)
}

// ExportCSV - exports test suite results as csv
func ExportCSV(rows [][]string) (string, error) {
	filePath, err := writeCSVToFile(rows)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func writeCSVToFile(rows [][]string) (string, error) {
	fileName := generateFileName(time.Now().UTC(), "csv")
	filePath := filepath.Join(config.Configuration.Export.Path, fileName)
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
