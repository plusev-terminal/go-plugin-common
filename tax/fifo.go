//go:build wasip1

package tax

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

//go:wasmimport extism:host/user fifo_reset
//go:noescape
func extismFifoReset(offset uint64) uint64

//go:wasmimport extism:host/user fifo_add_lot
//go:noescape
func extismFifoAddLot(offset uint64) uint64

//go:wasmimport extism:host/user fifo_dispose
//go:noescape
func extismFifoDispose(offset uint64) uint64

//go:wasmimport extism:host/user fifo_move_lots
//go:noescape
func extismFifoMoveLots(offset uint64) uint64

//go:wasmimport extism:host/user fifo_get_balance
//go:noescape
func extismFifoGetBalance(offset uint64) uint64

//go:wasmimport extism:host/user fifo_save_state
//go:noescape
func extismFifoSaveState(offset uint64) uint64

//go:wasmimport extism:host/user fifo_restore_state
//go:noescape
func extismFifoRestoreState(offset uint64) uint64

// FifoReset clears all FIFO engine state for a fresh run.
func FifoReset() error {
	mem := pdk.AllocateBytes([]byte("{}"))
	resOffset := extismFifoReset(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if !res.Result {
		return fmt.Errorf("fifo_reset error: %s", res.Error)
	}

	return nil
}

// FifoAddLot adds an acquisition lot to the FIFO queue.
func FifoAddLot(req FifoAddLotRequest) error {
	inputBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismFifoAddLot(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if !res.Result {
		return fmt.Errorf("fifo_add_lot error: %s", res.Error)
	}

	return nil
}

// FifoDispose consumes lots in FIFO order and returns match details.
func FifoDispose(req FifoDisposeRequest) (FifoDisposeResponse, error) {
	inputBytes, err := json.Marshal(req)
	if err != nil {
		return FifoDisposeResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismFifoDispose(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res FifoDisposeResponse
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return FifoDisposeResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !res.Result {
		return res, fmt.Errorf("fifo_dispose error: %s", res.Error)
	}

	return res, nil
}

// FifoMoveLots relocates lots between accounts, preserving acquisition timestamps.
// Partial moves are supported: if the source has fewer lots than requested,
// it moves what's available and reports the shortfall in the response.
func FifoMoveLots(req FifoMoveLotsRequest) (FifoMoveLotsResponse, error) {
	inputBytes, err := json.Marshal(req)
	if err != nil {
		return FifoMoveLotsResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismFifoMoveLots(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res FifoMoveLotsResponse
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return FifoMoveLotsResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if res.Error != "" {
		return res, fmt.Errorf("fifo_move_lots error: %s", res.Error)
	}

	return res, nil
}

// FifoGetBalance returns the remaining unconsumed amount for an account+asset pair.
func FifoGetBalance(req FifoGetBalanceRequest) (FifoGetBalanceResponse, error) {
	inputBytes, err := json.Marshal(req)
	if err != nil {
		return FifoGetBalanceResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismFifoGetBalance(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res FifoGetBalanceResponse
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return FifoGetBalanceResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !res.Result {
		return res, fmt.Errorf("fifo_get_balance error: %s", res.Error)
	}

	return res, nil
}

// FifoSaveState persists the current FIFO engine state under the given key.
func FifoSaveState(key string) error {
	inputBytes, err := json.Marshal(FifoStateRequest{Key: key})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismFifoSaveState(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if !res.Result {
		return fmt.Errorf("fifo_save_state error: %s", res.Error)
	}

	return nil
}

// FifoRestoreState restores FIFO engine state from a previous checkpoint.
// Returns true if the state was found and restored, false if not found.
func FifoRestoreState(key string) (bool, error) {
	inputBytes, err := json.Marshal(FifoStateRequest{Key: key})
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismFifoRestoreState(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res FifoStateResponse
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !res.Result {
		return false, fmt.Errorf("fifo_restore_state error: %s", res.Error)
	}

	return res.Found, nil
}
