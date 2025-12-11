package trading

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/plusev-terminal/go-plugin-common/utils"
)

type Unit string

func (u Unit) IsValid() error {
	switch u {
	case Hours, Minutes, Days, Weeks, Months, Years:
		return nil
	}
	return errors.New("unknown timeframe unit \"" + string(u) + "\"")
}

const (
	Hours   Unit = "h"
	Minutes Unit = "m"
	Days    Unit = "D"
	Weeks   Unit = "W"
	Months  Unit = "M"
	Years   Unit = "Y"
)

type Timeframe struct {
	ID       uint64         `json:"id" gorm:"primaryKey"`
	UserID   uint64         `json:"userId" gorm:"index"`
	Value    uint64         `json:"value" validate:"required"`
	Unit     Unit           `json:"unit" validate:"required"`
	Location *time.Location `json:"location" gorm:"-"` // Timezone for the timeframe
}

func NewTimeframe(val uint64, unit Unit, location ...*time.Location) Timeframe {
	loc := time.UTC
	if len(location) > 0 && location[0] != nil {
		loc = location[0]
	}
	return Timeframe{
		Value:    val,
		Unit:     unit,
		Location: loc,
	}
}

func (tf *Timeframe) String() string {
	return fmt.Sprintf("%d%s", tf.Value, tf.Unit)
}

func (tf *Timeframe) StringWithLocation() string {
	locStr := "UTC"
	if tf.Location != nil {
		locStr = tf.Location.String()
	}
	return fmt.Sprintf("%d%s:%s", tf.Value, tf.Unit, locStr)
}

func (tf Timeframe) ToMinutes(ref ...time.Time) int {
	mul := 1

	if tf.Unit == Hours {
		mul = 60
	}

	if tf.Unit == Days {
		mul = 60 * 24
	}

	if tf.Unit == Weeks {
		mul = 60 * 24 * 7
	}

	if tf.Unit == Months {
		if len(ref) > 0 {
			// Calculate exact minutes using the reference time
			startOfMonth := utils.StartOfMonth(ref[0])
			endTime := startOfMonth.AddDate(0, int(tf.Value), 0)
			return int(endTime.Sub(startOfMonth).Minutes())
		}
		mul = 60 * 24 * 30 // Approximate 30 days per month
	}

	if tf.Unit == Years {
		if len(ref) > 0 {
			// Calculate exact minutes using the reference time
			startOfYear := utils.StartOfYear(ref[0])
			endTime := startOfYear.AddDate(int(tf.Value), 0, 0)
			return int(endTime.Sub(startOfYear).Minutes())
		}
		mul = 60 * 24 * 365 // Approximate 365 days per year
	}

	return mul * int(tf.Value)
}

func (tf Timeframe) LowerThan(other Timeframe) bool {
	return tf.ToMinutes() < other.ToMinutes()
}

func (tf Timeframe) HigherThan(other Timeframe) bool {
	return tf.ToMinutes() > other.ToMinutes()
}

func (tf Timeframe) Equal(other Timeframe) bool {
	return tf.ToMinutes() == other.ToMinutes()
}

func (tf Timeframe) IsZero() bool {
	return tf.ToMinutes() == 0
}

// InLocation converts a time to the Timeframe's configured time zone
func (tf Timeframe) InLocation(tm time.Time) time.Time {
	if tf.Location == nil {
		return tm.In(time.UTC)
	}
	return tm.In(tf.Location)
}

func (tf Timeframe) IsValidCandleOpenTime(openTime time.Time) bool {
	// Convert to the configured time zone
	localTime := tf.InLocation(openTime)
	loc := tf.Location
	if loc == nil {
		loc = time.UTC
	}

	if tf.Unit == Minutes {
		day := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, loc)
		minuteOfDay := localTime.Sub(day).Minutes()
		return uint64(minuteOfDay)%tf.Value == 0
	}

	if tf.Unit == Hours {
		if localTime.Minute() == 0 && localTime.Hour()%int(tf.Value) == 0 {
			return true
		}
	}

	if tf.Unit == Days {
		if localTime.Minute() == 0 && localTime.Hour() == 0 && (localTime.YearDay()-1)%int(tf.Value) == 0 {
			return true
		}
	}

	if tf.Unit == Weeks {
		if localTime.Minute() == 0 && localTime.Hour() == 0 && localTime.Weekday() == time.Monday {
			return true
		}
	}

	if tf.Unit == Months {
		if localTime.Minute() == 0 && localTime.Hour() == 0 && localTime.Day() == 1 {
			return true
		}
	}

	if tf.Unit == Years {
		if localTime.Minute() == 0 && localTime.Hour() == 0 && localTime.Day() == 1 && localTime.Month() == time.January {
			return true
		}
	}

	return false
}

func (tf Timeframe) LastOpen(openTime time.Time) time.Time {
	localTime := tf.InLocation(openTime)

	if tf.Unit == Minutes {
		localTime = utils.StartOfMinute(localTime)
		delta := tf.ToMinutes() % 60
		distance := localTime.Minute() % delta * -1
		return localTime.Add(time.Duration(distance) * time.Minute)
	}

	if tf.Unit == Hours {
		localTime = utils.StartOfHour(localTime)
		distance := localTime.Hour() % int(tf.Value) * -1
		return localTime.Add(time.Duration(distance) * time.Hour)
	}

	if tf.Unit == Days {
		return utils.StartOfDay(localTime)
	}

	if tf.Unit == Weeks {
		return utils.StartOfWeek(localTime)
	}

	if tf.Unit == Months {
		return utils.StartOfMonth(localTime)
	}

	if tf.Unit == Years {
		return utils.StartOfYear(localTime)
	}

	return localTime
}

func (tf Timeframe) NextOpen(openTime time.Time) time.Time {
	lastOpen := tf.LastOpen(openTime)

	if lastOpen.Equal(openTime) {
		return openTime
	}

	// For Months and Years, use AddDate for proper calendar arithmetic
	if tf.Unit == Months {
		return lastOpen.AddDate(0, int(tf.Value), 0)
	}

	if tf.Unit == Years {
		return lastOpen.AddDate(int(tf.Value), 0, 0)
	}

	return lastOpen.Add(time.Duration(tf.ToMinutes()) * time.Minute)
}

func (tf Timeframe) CloseTime(openTime time.Time) time.Time {
	openTime = tf.LastOpen(openTime)

	// For Months and Years, use AddDate for proper calendar arithmetic
	if tf.Unit == Months {
		return openTime.AddDate(0, int(tf.Value), 0)
	}

	if tf.Unit == Years {
		return openTime.AddDate(int(tf.Value), 0, 0)
	}

	return openTime.Add(time.Duration(tf.ToMinutes()) * time.Minute)
}

func TimeframeFromString(str string) (Timeframe, error) {
	// Split the string into time frame and location parts (e.g., "4h:America/New_York" or "4h")
	parts := strings.Split(str, ":")
	if len(parts) < 1 {
		return Timeframe{}, errors.New("invalid timeframe format, expected 'valUnit[:location]'")
	}

	// Parse the time frame part (e.g., "4h")
	valUnit := parts[0]
	if len(valUnit) < 2 {
		return Timeframe{}, errors.New("invalid timeframe string")
	}

	valStr, err := strconv.ParseUint(valUnit[:len(valUnit)-1], 10, 64)
	if err != nil {
		return Timeframe{}, err
	}

	unit := Unit(valUnit[len(valUnit)-1:])
	if err := unit.IsValid(); err != nil {
		return Timeframe{}, err
	}

	// Parse the location (e.g., "America/New_York" or "UTC"), default to UTC if not provided
	var location *time.Location
	if len(parts) >= 2 && parts[1] != "" {
		if parts[1] == "UTC" {
			location = time.UTC
		} else {
			location, err = time.LoadLocation(parts[1])
			if err != nil {
				return Timeframe{}, errors.New("invalid time zone: " + err.Error())
			}
		}
	} else {
		location = time.UTC // Default to UTC if no location provided
	}

	return Timeframe{
		Value:    valStr,
		Unit:     unit,
		Location: location,
	}, nil
}
