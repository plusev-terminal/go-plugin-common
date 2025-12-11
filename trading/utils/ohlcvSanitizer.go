package utils

import (
	"fmt"
	"sort"

	tt "github.com/plusev-terminal/go-plugin-common/trading"
)

// OHLCVSanitizer processes OHLCV data batches to eliminate duplicates and fill gaps
type OHLCVSanitizer struct {
	timeframe   tt.Timeframe
	lastCandle  *tt.OHLCVRecord // Track the last processed candle to detect gaps
	firstCandle *tt.OHLCVRecord // Track the first processed candle for backward pagination
	initialized bool            // Whether we've processed at least one batch
}

// NewOHLCVSanitizer creates a new OHLCV sanitizer for the specified timeframe
func NewOHLCVSanitizer(timeframe tt.Timeframe) *OHLCVSanitizer {
	return &OHLCVSanitizer{
		timeframe:   timeframe,
		initialized: false,
	}
}

// SanitizeBatch processes a batch of OHLCV records, removing duplicates and filling gaps
func (s *OHLCVSanitizer) SanitizeBatch(batch []tt.OHLCVRecord) ([]tt.OHLCVRecord, error) {
	if len(batch) == 0 {
		return batch, nil
	}

	// Sort batch by opentime to ensure proper ordering
	sort.Slice(batch, func(i, j int) bool {
		return batch[i].OpenTime < batch[j].OpenTime
	})

	candleDurationSeconds := int64(s.timeframe.ToMinutes() * 60)
	result := make([]tt.OHLCVRecord, 0, len(batch))

	for i, candle := range batch {
		// 1. Internal Duplicate Check
		if i > 0 && candle.OpenTime == batch[i-1].OpenTime {
			continue
		}

		// 2. External Duplicate Check
		if s.initialized && s.firstCandle != nil && s.lastCandle != nil {
			if candle.OpenTime >= s.firstCandle.OpenTime && candle.OpenTime <= s.lastCandle.OpenTime {
				continue
			}
		}

		// 3. Gap Filling (Before the first valid candle of this batch)
		// Only fill gaps if we haven't added any candles to result yet (meaning this is the first new candle)
		// and we have a previous history to connect to.
		if len(result) == 0 && s.initialized && s.lastCandle != nil {
			if candle.OpenTime > s.lastCandle.OpenTime {
				nextTs := s.lastCandle.OpenTime + candleDurationSeconds
				for nextTs < candle.OpenTime {
					gap := tt.OHLCVRecord{
						OpenTime: nextTs,
						Open:     s.lastCandle.Close,
						High:     s.lastCandle.Close,
						Low:      s.lastCandle.Close,
						Close:    s.lastCandle.Close,
						Volume:   "0.00000000",
					}
					result = append(result, gap)
					nextTs += candleDurationSeconds
				}
			}
		}

		result = append(result, candle)
	}

	if len(result) == 0 {
		return []tt.OHLCVRecord{}, nil
	}

	// Update boundaries
	if !s.initialized || result[0].OpenTime < s.firstCandle.OpenTime {
		first := result[0]
		s.firstCandle = &first
	}
	if !s.initialized || result[len(result)-1].OpenTime > s.lastCandle.OpenTime {
		last := result[len(result)-1]
		s.lastCandle = &last
	}
	s.initialized = true

	return result, nil
}

// Reset clears the sanitizer state (useful for switching symbols/timeframes)
func (s *OHLCVSanitizer) Reset() {
	s.firstCandle = nil
	s.lastCandle = nil
	s.initialized = false
}

// GetLastCandle returns the last processed candle (useful for debugging)
func (s *OHLCVSanitizer) GetLastCandle() *tt.OHLCVRecord {
	if s.lastCandle == nil {
		return nil
	}
	// Return a copy to prevent external modification
	lastCopy := *s.lastCandle
	return &lastCopy
}

// SetTimeframe updates the timeframe (triggers reset)
func (s *OHLCVSanitizer) SetTimeframe(timeframe tt.Timeframe) {
	s.timeframe = timeframe
	s.Reset() // Reset state when timeframe changes
}

// ValidateBatch performs basic validation on OHLCV data
func (s *OHLCVSanitizer) ValidateBatch(batch []tt.OHLCVRecord) error {
	for i, record := range batch {
		if err := s.validateRecord(record); err != nil {
			return fmt.Errorf("invalid record at index %d: %w", i, err)
		}
	}
	return nil
}

// validateRecord checks if a single OHLCV record is valid
func (s *OHLCVSanitizer) validateRecord(record tt.OHLCVRecord) error {
	if record.OpenTime <= 0 {
		return fmt.Errorf("invalid opentime: %d", record.OpenTime)
	}

	// Parse prices to validate they're proper numbers
	open, err := parseFloat(record.Open)
	if err != nil {
		return fmt.Errorf("invalid open price: %s", record.Open)
	}

	high, err := parseFloat(record.High)
	if err != nil {
		return fmt.Errorf("invalid high price: %s", record.High)
	}

	low, err := parseFloat(record.Low)
	if err != nil {
		return fmt.Errorf("invalid low price: %s", record.Low)
	}

	close, err := parseFloat(record.Close)
	if err != nil {
		return fmt.Errorf("invalid close price: %s", record.Close)
	}

	// Validate OHLC relationships
	if high < low {
		return fmt.Errorf("high price (%.8f) cannot be less than low price (%.8f)", high, low)
	}

	if high < open || high < close {
		return fmt.Errorf("high price (%.8f) cannot be less than open (%.8f) or close (%.8f)", high, open, close)
	}

	if low > open || low > close {
		return fmt.Errorf("low price (%.8f) cannot be greater than open (%.8f) or close (%.8f)", low, open, close)
	}

	// Validate volume
	if _, err := parseFloat(record.Volume); err != nil {
		return fmt.Errorf("invalid volume: %s", record.Volume)
	}

	return nil
}

// Helper function to parse float from string (you might want to use a more robust parser)
func parseFloat(s string) (float64, error) {
	// This is a simplified parser - you might want to use strconv.ParseFloat
	// or a decimal library for better precision
	var f float64
	n, err := fmt.Sscanf(s, "%f", &f)
	if err != nil || n != 1 {
		return 0, fmt.Errorf("failed to parse float: %s", s)
	}
	return f, nil
}
