package helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/danilobml/lga/lga/internal/dtos"
)

func PrintStatusAnalysisResults(statusAnalysis *dtos.StatusAnalysisResponse, path string) {
	message := fmt.Sprintf("Total logs: %d, Total errors: %d, 4xx Errors: %d, 5xx Errors: %d\n", statusAnalysis.TotalLogs, statusAnalysis.ErrorsTotal, statusAnalysis.Errors400, statusAnalysis.Errors500)

	if path == "" {
		fmt.Println("Totals - " + message)
		return
	}

	fmt.Printf("Path: %s - %s", path, message)
}

func ParseDateTime(s string) (time.Time, error) {
	layouts := []string{
		"2006/01/02 - 15:04:05",
		"2006/01/02 15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"02/01/2006 15:04:05",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("no matching time layout for: " + s)
}
