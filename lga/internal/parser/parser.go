package parser

import (
	"bufio"
	"os"
	"regexp"

	"github.com/danilobml/lga/lga/internal/dtos"
	"github.com/danilobml/lga/lga/internal/helpers"
	"github.com/danilobml/lga/lga/internal/models"
)

func ParseFile(filePath string) (*dtos.ParserResponse, error) {
	response := dtos.ParserResponse{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	logs := []*models.Log{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		response.LinesRead++
		log, err := parseLine(scanner.Text())
		if err != nil {
			response.LinesSkipped++
			continue
		}
		response.LinesParsed++
		logs = append(logs, log)
	}

	response.Logs = logs

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &response, nil
}

func parseLine(line string) (*models.Log, error) {
	statusRe := regexp.MustCompile(`\|\s+(\d{3})\s+\|`)
	statusCode := statusRe.FindStringSubmatch(line)[1]

	pathRe := regexp.MustCompile(`"([^"]+)"\s*$`)
	path := pathRe.FindStringSubmatch(line)[1]

	dateTimeRe := regexp.MustCompile(`^(?:\[GIN\]\s+)?(\d{4}/\d{2}/\d{2}\s-\s\d{2}:\d{2}:\d{2})`)
	dateTime, err := helpers.ParseDateTime(dateTimeRe.FindStringSubmatch(line)[1])
	if err != nil {
		return nil, err
	}

	log := models.Log{
		StatusCode: statusCode,
		Path:       models.Path(path),
		DateTime:   dateTime,
	}

	return &log, nil
}
