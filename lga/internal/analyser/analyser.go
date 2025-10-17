package analyser

import (
	"fmt"
	"time"

	"github.com/danilobml/lga/lga/internal/dtos"
	"github.com/danilobml/lga/lga/internal/helpers"
	"github.com/danilobml/lga/lga/internal/models"
	"github.com/danilobml/lga/lga/internal/parser"
)

type Options struct {
	From  string
	To    string
	Paths bool
}

func AnalyseFileLogs(filePath string, options Options) error {
	logs, err := parser.ParseFile(filePath)
	if err != nil {
		return err
	}

	fmt.Println("options from", options.From)
	fmt.Println("options to", options.To)

	if options.From != "" || options.To != "" {
		from, _ := helpers.ParseDateTime(options.From)
		to, _ := helpers.ParseDateTime(options.To)
		logs = filterLogsPerPeriod(logs, to, from)
	}

	statusAnalysis := analyseStatus(logs)

	helpers.PrintStatusAnalysisResults(statusAnalysis, "")

	if options.Paths {
		pathAnalysis := analysePaths(logs)

		for path, statusAnalysis := range pathAnalysis.Results {
			helpers.PrintStatusAnalysisResults(statusAnalysis, string(path))
		}
	}

	return nil
}

func analyseStatus(logs []*models.Log) *dtos.StatusAnalysisResponse {
	response := dtos.StatusAnalysisResponse{}
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

func analysePaths(logs []*models.Log) *dtos.PathAnalysisResponse {
	type LogsPerPath map[models.Path][]*models.Log
	logsSorted := LogsPerPath{}

	for _, log := range logs {
		if logsSorted[log.Path] == nil {
			logsSorted[log.Path] = []*models.Log{log}
		} else {
			logsSorted[log.Path] = append(logsSorted[log.Path], log)
		}
	}

	response := dtos.PathAnalysisResponse{}
	response.Results = map[models.Path]*dtos.StatusAnalysisResponse{}

	for path, logs := range logsSorted {
		response.Results[path] = analyseStatus(logs)
	}

	return &response
}

func filterLogsPerPeriod(logs []*models.Log, to, from time.Time) []*models.Log {
	filteredLogs := []*models.Log{}

	fmt.Println("from", from)

	var realTo time.Time
	if to.IsZero() {
		realTo = time.Now()
	} else {
		realTo = to
	}

	fmt.Println("realTo", realTo)

	for _, log := range logs {
		if log.DateTime.After(from) && log.DateTime.Before(realTo) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	return filteredLogs
}
