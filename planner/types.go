package planner

import "time"

// ImportEvent represents an event to be imported
type ImportEvent struct {
	Title     string    `json:"title"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Notes     string    `json:"notes"`
	Timezone  string    `json:"timezone"`
	Tags      []string  `json:"tags"`
}
