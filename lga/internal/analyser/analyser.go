package analyser

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/danilobml/lga/lga/internal/models"
	"github.com/danilobml/lga/lga/internal/parser"
)

type Options struct {
	From  time.Time
	To    time.Time
	Paths bool
}

type StatusAnalysisResponse struct {
	TotalLogs   int
	ErrorsTotal int
	Errors400   int
	Errors500   int
}

type PathAnalysisResponse struct {
	Results map[string]StatusAnalysisResponse
}

type Log = models.Log

func AnalyseFileLogs(filePath string, options Options) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	logs := []*Log{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log, err := parser.ParseLine(scanner.Text())
		if err != nil {
			return err
		}
		logs = append(logs, log)
	}

	logAnalysis := analyseStatus(logs)

	// TODO: Extract to a helper:
	fmt.Printf("Totals: Total logs: %d, Total errors: %d, 4xx Errors: %d, 5xx Errors: %d\n", logAnalysis.TotalLogs, logAnalysis.ErrorsTotal, logAnalysis.Errors400, logAnalysis.Errors500)

	if options.Paths {
		pathAnalysis := analysePaths(logs)

		for path, statusAnalysis := range pathAnalysis.Results {
			fmt.Printf("Path: %s, Total logs: %d, Total errors: %d, 4xx Errors: %d, 5xx Errors: %d \n", path, statusAnalysis.TotalLogs, statusAnalysis.ErrorsTotal, statusAnalysis.Errors400, statusAnalysis.Errors500)
		}
	}

	// TODO: Filter per from and to, if set in options

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func analyseStatus(logs []*Log) *StatusAnalysisResponse {
	response := StatusAnalysisResponse{}
	errorsTotal := 0
	errors400 := 0
	errors500 := 0

	for _, log := range logs {
		if string(log.StatusCode[0]) == "4" || string(log.StatusCode[0]) == "5" {
			errorsTotal++
		}
		if string(log.StatusCode[0]) == "4" {
			errors400++
		}
		if string(log.StatusCode[0]) == "5" {
			errors500++
		}
	}

	response.TotalLogs = len(logs)
	response.ErrorsTotal = errorsTotal
	response.Errors400 = errors400
	response.Errors500 = errors500

	return &response
}

func analysePaths(logs []*Log) *PathAnalysisResponse {
	// TODO: Get paths from logs, filter by path

	response := PathAnalysisResponse{}
	response.Results = map[string]StatusAnalysisResponse{
		"test.com": {
			TotalLogs: 3,
			ErrorsTotal: 0,
			Errors400:   2,
			Errors500:   1,
		},
		"test2.com": {
			TotalLogs: 6,
			ErrorsTotal: 4,
			Errors400:   2,
			Errors500:   0,
		},
	}

	return &response
}
