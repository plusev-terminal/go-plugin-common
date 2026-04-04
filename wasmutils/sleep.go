package wasmutils

import "time"

// Import the sleep_ms host function
//
//go:wasmimport extism:host/user sleep_ms
func hostSleepMs(ms uint64) uint64

// Sleep pauses execution for the given duration by delegating to the host.
// WASM plugins can't reliably use time.Sleep, so this calls a host function instead.
func Sleep(d time.Duration) {
	ms := uint64(d.Milliseconds())
	if ms == 0 {
		return
	}
	hostSleepMs(ms)
}
