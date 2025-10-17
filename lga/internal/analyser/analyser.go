package analyser

import (
	"time"

	"github.com/danilobml/lga/lga/internal/dtos"
	"github.com/danilobml/lga/lga/internal/helpers"
	"github.com/danilobml/lga/lga/internal/models"
	"github.com/danilobml/lga/lga/internal/parser"
)

type Options struct {
	From  time.Time
	To    time.Time
	Paths bool
}

type Log = models.Log

func AnalyseFileLogs(filePath string, options Options) error {
	logs, err := parser.ParseFile(filePath)
	if err != nil {
		return err
	}

	statusAnalysis := analyseStatus(logs)

	helpers.PrintStatusAnalysisResults(statusAnalysis, "")

	if options.Paths {
		pathAnalysis := analysePaths(logs)

		for path, statusAnalysis := range pathAnalysis.Results {
			helpers.PrintStatusAnalysisResults(statusAnalysis, string(path))
		}
	}

	// TODO: Filter per from and to, if set in options

	return nil
}


func analyseStatus(logs []*Log) *dtos.StatusAnalysisResponse {
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

func analysePaths(logs []*Log) *dtos.PathAnalysisResponse {
	type LogsPerPath map[models.Path][]*Log
	logsSorted := LogsPerPath{}
	
	for _, log := range logs {
		if logsSorted[log.Path] == nil {
			logsSorted[log.Path] = []*Log{log}
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
