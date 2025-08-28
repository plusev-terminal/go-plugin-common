package types

import "time"

// ImportData defines the JSON structure for calendar import requests from a plugin.
type ImportData struct {
	Events []ImportEvent `json:"events"`
}

// ImportEvent represents an event to be imported
type ImportEvent struct {
	Title     string    `json:"title"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Notes     string    `json:"notes"`
	Timezone  string    `json:"timezone"`
	Tags      []string  `json:"tags"`
}

// ImportResult defines the JSON structure for calendar import responses to a plugin.
type ImportResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type ImportJob struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
