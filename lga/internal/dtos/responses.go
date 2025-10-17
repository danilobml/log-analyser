package dtos

import "github.com/danilobml/lga/lga/internal/models"

type StatusAnalysisResponse struct {
	TotalLogs   int
	ErrorsTotal int
	Errors400   int
	Errors500   int
}


type PathAnalysisResponse struct {
	Results map[models.Path]*StatusAnalysisResponse
}
