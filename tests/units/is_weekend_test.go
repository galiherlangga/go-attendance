package units

import (
	"testing"
	"time"

	"github.com/galiherlangga/go-attendance/pkg/utils"
)

func TestIsWeekend(t *testing.T) {
	tests := []struct {
		date     string
		expected bool
	}{
		{"2023-10-07", true},  // Saturday
		{"2023-10-08", true},  // Sunday
		{"2023-10-09", false}, // Monday
		{"2023-10-10", false}, // Tuesday
	}

	for _, test := range tests {
		date, err := time.Parse("2006-01-02", test.date)
		if err != nil {
			t.Fatalf("Failed to parse date %s: %v", test.date, err)
		}
		result := utils.IsWeekend(date)
		if result != test.expected {
			t.Errorf("IsWeekend(%s) = %v; want %v", test.date, result, test.expected)
		}
	}
}
