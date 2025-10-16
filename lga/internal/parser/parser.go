package parser

import (
	"regexp"

	"github.com/danilobml/lga/lga/internal/models"
)

type Log = models.Log

func ParseLine(line string) (*Log, error) {
	statusRe := regexp.MustCompile(`\|\s+(\d{3})\s+\|`)
	statusCode := statusRe.FindStringSubmatch(line)[1]

	log := Log{
		StatusCode: statusCode,
	}

	return &log, nil
}
