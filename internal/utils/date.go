package utils

import (
	"errors"
	"strings"
	"time"
)

func StartOfWeek(date time.Time) time.Time {
	weekday := int(date.Weekday())
	dayOffset := (weekday + 6) % 7
	d := date.AddDate(0, 0, -dayOffset)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, date.Location())
}

func EndOfWeek(date time.Time) time.Time {
	start := StartOfWeek(date)
	end := start.AddDate(0, 0, 6)
	return time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, date.Location())
}

func GetWeekDayIndex(date time.Time) int {
	weekday := int(date.Weekday())
	return (weekday + 6) % 7
}

func GetWeekParity(date time.Time) string {
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	diff := date.UTC().Sub(baseDate)
	diffWeeks := int(diff.Hours() / (24 * 7))
	if diffWeeks%2 == 0 {
		return "EVEN"
	}
	return "ODD"
}

func IsSameWeek(a, b time.Time) bool {
	startA := StartOfWeek(a)
	startB := StartOfWeek(b)
	return startA.Year() == startB.Year() &&
		startA.Month() == startB.Month() &&
		startA.Day() == startB.Day()
}

func IsSameDay(a, b time.Time) bool {
	return a.Year() == b.Year() &&
		a.Month() == b.Month() &&
		a.Day() == b.Day()
}

func IsDateInRange(date, start, end time.Time) bool {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	e := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	return !d.Before(s) && !d.After(e)
}

func DurationUntilEndOfDay(now time.Time) time.Duration {
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	dur := endOfDay.Sub(now)
	if dur <= 0 {
		return 1 * time.Minute
	}
	return dur
}

func ParseDate(val string) (time.Time, error) {
	val = strings.TrimSpace(val)
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, val); err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("invalid date format")
}
