package utils

import (
	"testing"
	"time"
)

func TestStartOfWeek(t *testing.T) {
	date := time.Date(2026, 4, 15, 14, 30, 0, 0, time.UTC)
	start := StartOfWeek(date)

	if start.Weekday() != time.Monday {
		t.Errorf("expected Monday, got %v", start.Weekday())
	}
	if start.Day() != 13 || start.Month() != 4 || start.Year() != 2026 {
		t.Errorf("expected 2026-04-13, got %v", start.Format("2006-01-02"))
	}
}

func TestGetWeekDayIndex(t *testing.T) {
	monday := time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC)
	sunday := time.Date(2026, 4, 19, 0, 0, 0, 0, time.UTC)

	if idx := GetWeekDayIndex(monday); idx != 0 {
		t.Errorf("expected 0 for Monday, got %d", idx)
	}
	if idx := GetWeekDayIndex(sunday); idx != 6 {
		t.Errorf("expected 6 for Sunday, got %d", idx)
	}
}

func TestGetWeekParity(t *testing.T) {
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	if parity := GetWeekParity(base); parity != "EVEN" {
		t.Errorf("expected EVEN, got %s", parity)
	}

	nextWeek := base.AddDate(0, 0, 7)
	if parity := GetWeekParity(nextWeek); parity != "ODD" {
		t.Errorf("expected ODD, got %s", parity)
	}
}
