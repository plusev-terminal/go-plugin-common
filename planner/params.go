package planner

import (
	"time"
)

// ImportParams contains parameters for the import command
type ImportParams struct {
	From time.Time `json:"from" validate:"required"`
	To   time.Time `json:"to" validate:"required"`
}
