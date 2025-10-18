package analyser

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/danilobml/lga/lga/internal/helpers"
	"github.com/danilobml/lga/lga/internal/models"
)

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	tt, err := helpers.ParseDateTime(s)
	if err != nil {
		t.Fatalf("failed to parse time %q: %v", s, err)
	}
	return tt
}

func mkLog(status, path, dt string) *models.Log {
	tt, _ := helpers.ParseDateTime(dt)
	return &models.Log{
		StatusCode: status,
		Path:       models.Path(path),
		DateTime:   tt,
	}
}

func TestAnalyseStatus(t *testing.T) {
	logs := []*models.Log{
		mkLog("200", "/ok", "2025/10/17 - 21:00:00"),
		mkLog("404", "/missing", "2025/10/17 - 21:01:00"),
		mkLog("500", "/boom", "2025/10/17 - 21:02:00"),
		mkLog("201", "/created", "2025/10/17 - 21:03:00"),
		mkLog("418", "/teapot", "2025/10/17 - 21:04:00"),
	}

	res := analyseStatus(logs)

	if res.TotalLogs != 5 {
		t.Errorf("TotalLogs: want 5, got %d", res.TotalLogs)
	}
	if res.ErrorsTotal != 3 {
		t.Errorf("ErrorsTotal: want 3, got %d", res.ErrorsTotal)
	}
	if res.Errors400 != 2 {
		t.Errorf("Errors400: want 2, got %d", res.Errors400)
	}
	if res.Errors500 != 1 {
		t.Errorf("Errors500: want 1, got %d", res.Errors500)
	}
}

func TestAnalysePaths(t *testing.T) {
	logs := []*models.Log{
		mkLog("200", "/users", "2025/10/17 - 21:00:00"),
		mkLog("404", "/users", "2025/10/17 - 21:01:00"),
		mkLog("500", "/orders", "2025/10/17 - 21:02:00"),
		mkLog("201", "/orders", "2025/10/17 - 21:03:00"),
		mkLog("418", "/orders", "2025/10/17 - 21:04:00"),
	}

	res := analysePaths(logs)
	if res == nil || res.Results == nil {
		t.Fatalf("analysePaths returned nil or empty results")
	}

	users := res.Results[models.Path("/users")]
	if users == nil {
		t.Fatalf("missing /users analysis")
	}
	if users.TotalLogs != 2 || users.ErrorsTotal != 1 || users.Errors400 != 1 || users.Errors500 != 0 {
		t.Errorf("/users analysis wrong: %+v", users)
	}

	orders := res.Results[models.Path("/orders")]
	if orders == nil {
		t.Fatalf("missing /orders analysis")
	}
	if orders.TotalLogs != 3 || orders.ErrorsTotal != 2 || orders.Errors400 != 1 || orders.Errors500 != 1 {
		t.Errorf("/orders analysis wrong: %+v", orders)
	}
}

func TestFilterLogsPerPeriod(t *testing.T) {
	logA := mkLog("200", "/a", "2025/10/17 - 21:00:00")
	logB := mkLog("404", "/b", "2025/10/17 - 22:00:00")
	logC := mkLog("500", "/c", "2025/10/17 - 23:00:00")

	from := mustTime(t, "2025/10/17 - 21:00:00")
	to := mustTime(t, "2025/10/17 - 23:00:00")

	got := filterLogsPerPeriod([]*models.Log{logA, logB, logC}, to, from)
	if len(got) != 1 || got[0] != logB {
		t.Errorf("expected only logB within (from,to) exclusive, got %v", len(got))
	}

	zeroTo := time.Time{}
	got2 := filterLogsPerPeriod([]*models.Log{logA, logB}, zeroTo, mustTime(t, "2025/10/16 - 00:00:00"))
	if len(got2) != 2 {
		t.Errorf("expected 2 logs when to is zero and from is yesterday, got %d", len(got2))
	}
}

func TestAnalyseFileLogs_Integration(t *testing.T) {
	lines := []string{
		`[GIN] 2025/10/17 - 21:00:00 | 200 | "GET /users"`,
		`[GIN] 2025/10/17 - 21:05:00 | 404 | "GET /users/123"`,
		`[GIN] 2025/10/17 - 21:10:00 | 500 | "POST /orders"`,
		`this is not a valid gin log line`,
		`[GIN] 2025/10/17 - 22:00:00 | 201 | "POST /users"`,
	}

	dir := t.TempDir()
	fp := filepath.Join(dir, "app.log")
	if err := os.WriteFile(fp, []byte(joinLines(lines)), 0o644); err != nil {
		t.Fatalf("failed to write temp log file: %v", err)
	}

	prevStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	opts := Options{
		From:  "2025/10/17 - 20:59:59",
		To:    "2025/10/17 - 22:00:00", 
		Paths: true,
	}

	err := AnalyseFileLogs(fp, opts)

	_ = w.Close()
	os.Stdout = prevStdout
	_, _ = io.ReadAll(r)

	if err != nil {
		t.Fatalf("AnalyseFileLogs returned error: %v", err)
	}
}

func joinLines(lines []string) string {
	var b bytes.Buffer
	for i, s := range lines {
		b.WriteString(s)
		if i < len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
