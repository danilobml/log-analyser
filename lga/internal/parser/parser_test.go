package parser

import (
	"testing"
	"time"

	"github.com/danilobml/lga/lga/internal/models"
)

func TestParseLine_Success(t *testing.T) {
	line := `[GIN] 2025/10/17 - 21:04:05 | 200 | "GET /api/users"`
	
	log, err := parseLine(line)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedStatus := "200"
	expectedPath := models.Path("GET /api/users")
	expectedTime, _ := time.Parse("2006/01/02 - 15:04:05", "2025/10/17 - 21:04:05")

	if log.StatusCode != expectedStatus {
		t.Errorf("expected StatusCode %s, got %s", expectedStatus, log.StatusCode)
	}

	if log.Path != expectedPath {
		t.Errorf("expected Path %s, got %s", expectedPath, log.Path)
	}

	if !log.DateTime.Equal(expectedTime) {
		t.Errorf("expected DateTime %v, got %v", expectedTime, log.DateTime)
	}
}

func TestParseLine_InvalidFormat(t *testing.T) {
	line := `[GIN] 2025/10/17 - 21:04:05 something invalid`

	_, err := parseLine(line)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseLine_InvalidDate(t *testing.T) {
	line := `[GIN] 2025-10-17 21:04:05 | 404 | "GET /api/test"`

	_, err := parseLine(line)
	if err == nil {
		t.Fatal("expected error for invalid date format, got nil")
	}
}
