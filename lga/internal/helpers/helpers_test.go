package helpers

import "testing"

func TestParseDateTime(t *testing.T) {
	validCases := []string{
		"2025/10/18 - 01:44:05",
		"2025/10/18 01:44:05",
		"2025-10-18 01:44:05",
		"2025-10-18T01:44:05",
		"2025-10-18T01:44:05+02:00",
		"18/10/2025 01:44:05",
		"18/10/2025",
		"10/18/2025",
	}

	for _, input := range validCases {
		t.Run("valid_"+input, func(t *testing.T) {
			got, err := ParseDateTime(input)
			if err != nil {
				t.Fatalf("expected no error for %q, got %v", input, err)
			}

			if got.IsZero() {
				t.Errorf("expected non-zero time for %q", input)
			}
		})
	}

	invalidCases := []string{
		"not a date",
		"2025-18-10",
		"",
	}

	for _, input := range invalidCases {
		t.Run("invalid_"+input, func(t *testing.T) {
			_, err := ParseDateTime(input)
			if err == nil {
				t.Fatalf("expected error for %q, got nil", input)
			}
		})
	}
}
