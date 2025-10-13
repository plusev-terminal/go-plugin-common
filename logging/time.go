package logging

import (
	"encoding/json"
	"time"

	"github.com/extism/go-pdk"
)

// Import the time_now host function
//
//go:wasmimport extism:host/user time_now
func hostTimeNow(_ uint64) uint64

// Now returns the current time from the host
// This is necessary because WASM plugins don't have reliable access to system time
func Now() (time.Time, error) {
	// Call the host function to get current time (pass 0 as dummy parameter)
	offset := hostTimeNow(0)
	if offset == 0 {
		// Host function failed, return zero time
		return time.Time{}, nil
	}

	// Read the time data from memory
	timeMem := pdk.FindMemory(offset)
	timeBytes := timeMem.ReadBytes()

	// Unmarshal the JSON time
	var t time.Time
	if err := json.Unmarshal(timeBytes, &t); err != nil {
		return time.Time{}, err
	}

	return t, nil
}
