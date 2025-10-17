package trading

import "fmt"

type Unit string

const (
	Hours   Unit = "h"
	Minutes Unit = "m"
	Days    Unit = "D"
	Weeks   Unit = "W"
	Months  Unit = "M"
	Years   Unit = "Y"
)

// Timeframe represents a supported timeframe/interval
type Timeframe struct {
	Value uint64 `json:"value"`
	Unit  Unit   `json:"unit"`
}

func (tf *Timeframe) String() string {
	return fmt.Sprintf("%d%s", tf.Value, tf.Unit)
}
