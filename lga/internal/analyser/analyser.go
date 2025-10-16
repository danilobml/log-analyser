package analyser

import (
	"bufio"
	"fmt"
	"os"

	"github.com/danilobml/lga/lga/internal/models"
	"github.com/danilobml/lga/lga/internal/parser"
)

type StatusAnalysisResponse struct {
	ErrorsTotal int
	Errors400   int
	Errors500   int
}

type Log = models.Log

func AnalyseFileLogs(filePath string) error {
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

	fmt.Printf("Total errors: %d, 4xx Errors: %d, 5xx Errors: %d", logAnalysis.ErrorsTotal, logAnalysis.Errors400, logAnalysis.Errors500)

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

	response.ErrorsTotal = errorsTotal
	response.Errors400 = errors400
	response.Errors500 = errors500

	return &response
}
