//go:build wasip1

package tax

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

//go:wasmimport extism:host/user submit_trade
//go:noescape
func extismSubmitTrade(offset uint64) uint64

//go:wasmimport extism:host/user submit_transfer
//go:noescape
func extismSubmitTransfer(offset uint64) uint64

// SubmitTrade sends a trade record to the host for storage.
func SubmitTrade(trade PluginTrade) error {
	inputBytes, err := json.Marshal(trade)
	if err != nil {
		return fmt.Errorf("failed to marshal trade: %w", err)
	}
	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismSubmitTrade(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()
	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal submit result: %w", err)
	}
	if !res.Result {
		return fmt.Errorf("trade submit error: %s", res.Error)
	}
	return nil
}

// SubmitTransfer sends a transfer record to the host for storage.
func SubmitTransfer(transfer PluginTransfer) error {
	inputBytes, err := json.Marshal(transfer)
	if err != nil {
		return fmt.Errorf("failed to marshal transfer: %w", err)
	}
	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismSubmitTransfer(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()
	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal submit result: %w", err)
	}
	if !res.Result {
		return fmt.Errorf("transfer submit error: %s", res.Error)
	}
	return nil
}
