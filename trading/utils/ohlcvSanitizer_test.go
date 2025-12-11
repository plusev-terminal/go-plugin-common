package utils

import (
	"testing"

	tt "github.com/plusev-terminal/go-plugin-common/trading"
)

func TestOHLCVSanitizer_RemoveDuplicates(t *testing.T) {
	timeframe, _ := tt.TimeframeFromString("5m")
	sanitizer := NewOHLCVSanitizer(timeframe)

	// Test batch with duplicate first candle
	batch1 := []tt.OHLCVRecord{
		{Timestamp: 1000, Open: "100.0", High: "101.0", Low: "99.0", Close: "100.5", Volume: "1000"},
		{Timestamp: 1300, Open: "100.5", High: "102.0", Low: "100.0", Close: "101.0", Volume: "2000"},
	}

	batch2 := []tt.OHLCVRecord{
		{Timestamp: 1300, Open: "100.5", High: "102.0", Low: "100.0", Close: "101.0", Volume: "2000"}, // Duplicate
		{Timestamp: 1600, Open: "101.0", High: "103.0", Low: "101.0", Close: "102.0", Volume: "1500"},
	}

	// Process first batch
	result1, err := sanitizer.SanitizeBatch(batch1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result1) != 2 {
		t.Fatalf("Expected 2 records in first batch, got %d", len(result1))
	}

	// Process second batch (should remove duplicate)
	result2, err := sanitizer.SanitizeBatch(batch2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result2) != 1 {
		t.Fatalf("Expected 1 record in second batch after duplicate removal, got %d", len(result2))
	}

	if result2[0].Timestamp != 1600 {
		t.Fatalf("Expected timestamp 1600, got %d", result2[0].Timestamp)
	}
}

func TestOHLCVSanitizer_FillGaps(t *testing.T) {
	timeframe, _ := tt.TimeframeFromString("5m") // 5 minutes = 300 seconds
	sanitizer := NewOHLCVSanitizer(timeframe)

	// First batch
	batch1 := []tt.OHLCVRecord{
		{Timestamp: 1000, Open: "100.0", High: "101.0", Low: "99.0", Close: "100.5", Volume: "1000"},
	}

	// Second batch with gap (should be at 1300, but starts at 1900 - missing 2 candles)
	batch2 := []tt.OHLCVRecord{
		{Timestamp: 1900, Open: "102.0", High: "103.0", Low: "101.5", Close: "102.5", Volume: "1500"},
	}

	// Process first batch
	_, err := sanitizer.SanitizeBatch(batch1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Process second batch (should fill gaps)
	result2, err := sanitizer.SanitizeBatch(batch2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have filled 2 gaps + 1 real candle = 3 records
	if len(result2) != 3 {
		t.Fatalf("Expected 3 records (2 gap fills + 1 real), got %d", len(result2))
	}

	// Check gap fill candles
	expectedTimestamps := []int64{1300, 1600, 1900}
	for i, expected := range expectedTimestamps {
		if result2[i].Timestamp != expected {
			t.Fatalf("Expected timestamp %d at index %d, got %d", expected, i, result2[i].Timestamp)
		}
	}

	// Gap fill candles should have volume 0 and use previous close price
	if result2[0].Volume != "0.00000000" {
		t.Fatalf("Expected gap fill volume to be 0, got %s", result2[0].Volume)
	}

	if result2[0].Close != "100.5" { // Previous candle's close
		t.Fatalf("Expected gap fill close to be 100.5, got %s", result2[0].Close)
	}
}

func TestOHLCVSanitizer_OverlapAndGaps(t *testing.T) {
	timeframe, _ := tt.TimeframeFromString("5m") // 300s
	sanitizer := NewOHLCVSanitizer(timeframe)

	// Batch 1: 1000, 1300
	batch1 := []tt.OHLCVRecord{
		{Timestamp: 1000, Open: "100", High: "100", Low: "100", Close: "100", Volume: "100"},
		{Timestamp: 1300, Open: "100", High: "100", Low: "100", Close: "100", Volume: "100"},
	}

	// Batch 2: 1300 (overlap), 1900 (gap of 1600)
	batch2 := []tt.OHLCVRecord{
		{Timestamp: 1300, Open: "100", High: "100", Low: "100", Close: "100", Volume: "100"},
		{Timestamp: 1900, Open: "100", High: "100", Low: "100", Close: "100", Volume: "100"},
	}

	_, _ = sanitizer.SanitizeBatch(batch1)
	result, _ := sanitizer.SanitizeBatch(batch2)

	// Expected:
	// 1300 is skipped (duplicate)
	// Gap fill for 1600 is inserted
	// 1900 is appended
	// Total result length: 2 (1600, 1900)

	if len(result) != 2 {
		t.Fatalf("Expected 2 records (1 gap fill + 1 new), got %d", len(result))
	}

	if result[0].Timestamp != 1600 {
		t.Errorf("Expected gap fill at 1600, got %d", result[0].Timestamp)
	}
	if result[1].Timestamp != 1900 {
		t.Errorf("Expected new candle at 1900, got %d", result[1].Timestamp)
	}
}

func TestOHLCVSanitizer_BackwardFetch(t *testing.T) {
	timeframe, _ := tt.TimeframeFromString("5m")
	sanitizer := NewOHLCVSanitizer(timeframe)

	// Fetching backwards: newer data comes first
	batch1 := []tt.OHLCVRecord{
		{Timestamp: 2000, Open: "100", High: "100", Low: "100", Close: "100", Volume: "100"},
	}
	batch2 := []tt.OHLCVRecord{
		{Timestamp: 1700, Open: "100", High: "100", Low: "100", Close: "100", Volume: "100"},
	}

	_, _ = sanitizer.SanitizeBatch(batch1)
	result, _ := sanitizer.SanitizeBatch(batch2)

	// Should accept older data without gap filling (gap filling is only forward)
	if len(result) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(result))
	}
	if result[0].Timestamp != 1700 {
		t.Errorf("Expected timestamp 1700, got %d", result[0].Timestamp)
	}
}

func TestOHLCVSanitizer_Validation(t *testing.T) {
	timeframe, _ := tt.TimeframeFromString("1m")
	sanitizer := NewOHLCVSanitizer(timeframe)

	// Test invalid OHLC relationships
	invalidBatch := []tt.OHLCVRecord{
		{Timestamp: 1000, Open: "100.0", High: "99.0", Low: "101.0", Close: "100.5", Volume: "1000"}, // High < Low
	}

	err := sanitizer.ValidateBatch(invalidBatch)
	if err == nil {
		t.Fatalf("Expected validation error for invalid OHLC relationships")
	}

	// Test invalid timestamp
	invalidBatch2 := []tt.OHLCVRecord{
		{Timestamp: 0, Open: "100.0", High: "101.0", Low: "99.0", Close: "100.5", Volume: "1000"},
	}

	err = sanitizer.ValidateBatch(invalidBatch2)
	if err == nil {
		t.Fatalf("Expected validation error for invalid timestamp")
	}
}

func TestOHLCVSanitizer_Reset(t *testing.T) {
	timeframe, _ := tt.TimeframeFromString("1h")
	sanitizer := NewOHLCVSanitizer(timeframe)

	// Process a batch
	batch := []tt.OHLCVRecord{
		{Timestamp: 1000, Open: "100.0", High: "101.0", Low: "99.0", Close: "100.5", Volume: "1000"},
	}

	_, err := sanitizer.SanitizeBatch(batch)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have last candle set
	if sanitizer.GetLastCandle() == nil {
		t.Fatalf("Expected last candle to be set")
	}

	// Reset
	sanitizer.Reset()

	// Should be cleared
	if sanitizer.GetLastCandle() != nil {
		t.Fatalf("Expected last candle to be nil after reset")
	}

	if sanitizer.initialized {
		t.Fatalf("Expected initialized to be false after reset")
	}
}
