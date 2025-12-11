package utils

import (
	"sort"
	"strings"
	"time"
)

func StartOfMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

func StartOfHour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
}

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func StartOfWeek(t time.Time) time.Time {
	// Get the weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	weekday := t.Weekday()

	// Calculate days to subtract to get to Monday
	// If it's Sunday (0), we need to go back 6 days
	// If it's Monday (1), we need to go back 0 days
	// If it's Tuesday (2), we need to go back 1 day, etc.
	daysToSubtract := int(weekday)
	if weekday == time.Sunday {
		daysToSubtract = 6
	} else {
		daysToSubtract = int(weekday) - 1
	}

	// Get the start of the current day and subtract the required days
	startOfDay := StartOfDay(t)
	return startOfDay.AddDate(0, 0, -daysToSubtract)
}

func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func StartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, time.January, 1, 0, 0, 0, 0, t.Location())
}

var humanDateFormatTokenMap = map[string]string{
	"dddd": "Monday",
	"ddd":  "Mon",
	"dd":   "02",
	"d":    "2",
	"MMMM": "January",
	"MMM":  "Jan",
	"MM":   "01",
	"M":    "1",
	"YYYY": "2006",
	"YY":   "06",
	"HH":   "15",
	"H":    "15",
	"hh":   "03",
	"h":    "3",
	"mm":   "04",
	"m":    "4",
	"ss":   "05",
	"s":    "5",
	"AMPM": "PM",
	"ampm": "pm",
	"AP":   "PM",
	"ap":   "pm",
	"Z":    "Z07:00",
	"z":    "Z07:00",
}

func HumanDateFormatToGoFormat(format string) string {
	// Check for predefined format names first
	switch strings.ToLower(format) {
	case "iso", "iso8601":
		return "2006-01-02T15:04:05Z07:00"
	case "rfc3339":
		return time.RFC3339
	case "short":
		return "2006-01-02"
	case "time":
		return "15:04:05"
	case "datetime":
		return "2006-01-02 15:04:05"
	}

	// Get tokens sorted by length (longest first) to handle overlapping patterns
	var tokens []string
	for token := range humanDateFormatTokenMap {
		tokens = append(tokens, token)
	}
	sort.Slice(tokens, func(i, j int) bool {
		return len(tokens[i]) > len(tokens[j])
	})

	// Create a new result string by scanning through the input format
	var result strings.Builder
	remaining := format

	for len(remaining) > 0 {
		foundMatch := false

		// Try to match each token at the current position
		for _, token := range tokens {
			// Check for case-insensitive match at the start of remaining string
			if len(remaining) >= len(token) &&
				strings.EqualFold(remaining[:len(token)], token) {
				// Found a match, add the Go format string
				result.WriteString(humanDateFormatTokenMap[token])
				// Remove the matched portion from the remaining string
				remaining = remaining[len(token):]
				foundMatch = true
				break
			}
		}

		// If no token matches at current position, copy the character as-is
		if !foundMatch {
			result.WriteByte(remaining[0])
			remaining = remaining[1:]
		}
	}

	return result.String()
}
