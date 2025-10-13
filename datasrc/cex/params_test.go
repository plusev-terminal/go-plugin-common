package cex_test

import (
	"testing"
	"time"

	"github.com/plusev-terminal/go-plugin-common/datasrc/cex"
)

func TestParseOHLCVStreamParams(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		want    *cex.OHLCVStreamParams
		wantErr bool
	}{
		{
			name: "valid params",
			input: map[string]any{
				"symbol":   "BTC/USDT",
				"interval": "1m",
			},
			want: &cex.OHLCVStreamParams{
				Symbol:    "BTC/USDT",
				Timeframe: "1m",
			},
			wantErr: false,
		},
		{
			name: "missing symbol",
			input: map[string]any{
				"interval": "1m",
			},
			wantErr: true,
		},
		{
			name: "missing interval",
			input: map[string]any{
				"symbol": "BTC/USDT",
			},
			wantErr: true,
		},
		{
			name: "invalid symbol type",
			input: map[string]any{
				"symbol":   123,
				"interval": "1m",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cex.ParseOHLCVStreamParams(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOHLCVStreamParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Symbol != tt.want.Symbol || got.Interval != tt.want.Timeframe {
					t.Errorf("ParseOHLCVStreamParams() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestParseGetOHLCVParams(t *testing.T) {
	now := time.Now()
	nowStr := now.Format(time.RFC3339)
	nowMillis := now.UnixMilli()

	tests := []struct {
		name    string
		input   map[string]any
		want    *cex.GetOHLCVParams
		wantErr bool
	}{
		{
			name: "minimal params",
			input: map[string]any{
				"symbol":    "BTC/USDT",
				"timeframe": "1h",
			},
			want: &cex.GetOHLCVParams{
				Symbol:    "BTC/USDT",
				Timeframe: "1h",
			},
			wantErr: false,
		},
		{
			name: "with time string",
			input: map[string]any{
				"symbol":    "BTC/USDT",
				"timeframe": "1h",
				"startTime": nowStr,
			},
			want: &cex.GetOHLCVParams{
				Symbol:    "BTC/USDT",
				Timeframe: "1h",
				StartTime: &now,
			},
			wantErr: false,
		},
		{
			name: "with unix timestamp",
			input: map[string]any{
				"symbol":    "BTC/USDT",
				"timeframe": "1h",
				"startTime": float64(nowMillis), // JSON numbers are float64
			},
			want: &cex.GetOHLCVParams{
				Symbol:    "BTC/USDT",
				Timeframe: "1h",
				StartTime: &now,
			},
			wantErr: false,
		},
		{
			name: "with limit",
			input: map[string]any{
				"symbol":    "BTC/USDT",
				"timeframe": "1h",
				"limit":     float64(100), // JSON numbers are float64
			},
			want: &cex.GetOHLCVParams{
				Symbol:    "BTC/USDT",
				Timeframe: "1h",
				Limit:     100,
			},
			wantErr: false,
		},
		{
			name: "missing symbol",
			input: map[string]any{
				"timeframe": "1h",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cex.ParseGetOHLCVParams(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGetOHLCVParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Symbol != tt.want.Symbol || got.Timeframe != tt.want.Timeframe {
					t.Errorf("ParseGetOHLCVParams() = %+v, want %+v", got, tt.want)
				}
				if got.Limit != tt.want.Limit {
					t.Errorf("ParseGetOHLCVParams() limit = %d, want %d", got.Limit, tt.want.Limit)
				}
			}
		})
	}
}

func TestOHLCVStreamParamsToMap(t *testing.T) {
	params := &cex.OHLCVStreamParams{
		Symbol:    "BTC/USDT",
		Timeframe: "1m",
	}

	result := params.ToMap()

	if result["symbol"] != "BTC/USDT" {
		t.Errorf("ToMap() symbol = %v, want BTC/USDT", result["symbol"])
	}
	if result["interval"] != "1m" {
		t.Errorf("ToMap() interval = %v, want 1m", result["interval"])
	}
}

func TestGetOHLCVParamsToMap(t *testing.T) {
	now := time.Now()
	params := &cex.GetOHLCVParams{
		Symbol:    "BTC/USDT",
		Timeframe: "1h",
		StartTime: &now,
		Limit:     100,
	}

	result := params.ToMap()

	if result["symbol"] != "BTC/USDT" {
		t.Errorf("ToMap() symbol = %v, want BTC/USDT", result["symbol"])
	}
	if result["timeframe"] != "1h" {
		t.Errorf("ToMap() timeframe = %v, want 1h", result["timeframe"])
	}
	if result["limit"] != 100 {
		t.Errorf("ToMap() limit = %v, want 100", result["limit"])
	}
	if _, ok := result["startTime"]; !ok {
		t.Error("ToMap() missing startTime")
	}
}
