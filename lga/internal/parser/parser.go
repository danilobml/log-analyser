package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/danilobml/lga/lga/internal/dtos"
	"github.com/danilobml/lga/lga/internal/helpers"
	"github.com/danilobml/lga/lga/internal/models"
)

var (
	statusRe   = regexp.MustCompile(`\|\s+(\d{3})\s+\|`)
	pathRe     = regexp.MustCompile(`"([^"]+)"\s*$`)
	dateTimeRe = regexp.MustCompile(`^(?:\[GIN\]\s+)?(\d{4}/\d{2}/\d{2}\s-\s\d{2}:\d{2}:\d{2})`)
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
	statusMatch := statusRe.FindStringSubmatch(line)
	if len(statusMatch) < 2 {
		return nil, fmt.Errorf("parseLine: status code not found")
	}
	statusCode := statusMatch[1]

	pathMatch := pathRe.FindStringSubmatch(line)
	if len(pathMatch) < 2 {
		return nil, fmt.Errorf("parseLine: path not found")
	}
	path := pathMatch[1]

	dateTimeMatch := dateTimeRe.FindStringSubmatch(line)
	if len(dateTimeMatch) < 2 {
		return nil, fmt.Errorf("parseLine: datetime not found")
	}
	dateTime, err := helpers.ParseDateTime(dateTimeMatch[1])
	if err != nil {
		return nil, fmt.Errorf("parseLine: %w", err)
	}

	log := models.Log{
		StatusCode: statusCode,
		Path:       models.Path(path),
		DateTime:   dateTime,
	}

	return &log, nil
}
