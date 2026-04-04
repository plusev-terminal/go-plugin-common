//go:build wasip1

package tax

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

//go:wasmimport extism:host/user get_records
//go:noescape
func extismGetRecords(offset uint64) uint64

//go:wasmimport extism:host/user get_settings
//go:noescape
func extismGetSettings(offset uint64) uint64

//go:wasmimport extism:host/user submit_report_entry
//go:noescape
func extismSubmitReportEntry(offset uint64) uint64

//go:wasmimport extism:host/user report_progress
//go:noescape
func extismReportProgress(offset uint64) uint64

//go:wasmimport extism:host/user submit_report_summary
//go:noescape
func extismSubmitReportSummary(offset uint64) uint64

//go:wasmimport extism:host/user kv_put
//go:noescape
func extismKVPut(offset uint64) uint64

//go:wasmimport extism:host/user kv_get
//go:noescape
func extismKVGet(offset uint64) uint64

//go:wasmimport extism:host/user kv_delete
//go:noescape
func extismKVDelete(offset uint64) uint64

//go:wasmimport extism:host/user kv_list
//go:noescape
func extismKVList(offset uint64) uint64

// GetRecords fetches paginated history records from the host.
func GetRecords(req PluginGetRecordsRequest) (PluginGetRecordsResponse, error) {
	inputBytes, err := json.Marshal(req)
	if err != nil {
		return PluginGetRecordsResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismGetRecords(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginGetRecordsResponse
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return PluginGetRecordsResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !res.Result {
		return res, fmt.Errorf("get_records error: %s", res.Error)
	}

	return res, nil
}

// GetSettings retrieves user tax settings from the host.
func GetSettings() (PluginSettingsResponse, error) {
	mem := pdk.AllocateBytes([]byte("{}"))
	resOffset := extismGetSettings(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginSettingsResponse
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return PluginSettingsResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !res.Result {
		return res, fmt.Errorf("get_settings error: %s", res.Error)
	}

	return res, nil
}

// SubmitReportEntry sends a report entry to the host.
func SubmitReportEntry(entry PluginReportEntry) error {
	inputBytes, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismSubmitReportEntry(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if !res.Result {
		return fmt.Errorf("submit_report_entry error: %s", res.Error)
	}

	return nil
}

// ReportProgress sends progress updates to the host for UI display.
func ReportProgress(progress PluginReportProgress) error {
	inputBytes, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismReportProgress(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if !res.Result {
		return fmt.Errorf("report_progress error: %s", res.Error)
	}

	return nil
}

// SubmitReportSummary sends the report summary rows to the host for storage.
func SubmitReportSummary(rows []ReportSummaryRow) error {
	summary := PluginReportSummary{Rows: rows}

	inputBytes, err := json.Marshal(summary)
	if err != nil {
		return fmt.Errorf("failed to marshal summary: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismSubmitReportSummary(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if !res.Result {
		return fmt.Errorf("submit_report_summary error: %s", res.Error)
	}

	return nil
}

// KVPut stores a value in the plugin's key-value store.
func KVPut(namespace, key string, value []byte) error {
	inputBytes, err := json.Marshal(PluginKVPutRequest{
		Namespace: namespace,
		Key:       key,
		Value:     value,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismKVPut(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if !res.Result {
		return fmt.Errorf("kv_put error: %s", res.Error)
	}

	return nil
}

// KVGet retrieves a value from the plugin's key-value store.
// Returns the value, whether the key was found, and any error.
func KVGet(namespace, key string) ([]byte, bool, error) {
	inputBytes, err := json.Marshal(PluginKVGetRequest{
		Namespace: namespace,
		Key:       key,
	})
	if err != nil {
		return nil, false, fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismKVGet(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginKVGetResponse
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !res.Result {
		return nil, false, fmt.Errorf("kv_get error: %s", res.Error)
	}

	return res.Value, res.Found, nil
}

// KVDelete removes a key from the plugin's key-value store.
func KVDelete(namespace, key string) error {
	inputBytes, err := json.Marshal(PluginKVDeleteRequest{
		Namespace: namespace,
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismKVDelete(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginSubmitResult
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	if !res.Result {
		return fmt.Errorf("kv_delete error: %s", res.Error)
	}

	return nil
}

// KVList lists keys in a namespace with an optional prefix filter.
func KVList(namespace, prefix string) ([]string, error) {
	inputBytes, err := json.Marshal(PluginKVListRequest{
		Namespace: namespace,
		Prefix:    prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	mem := pdk.AllocateBytes(inputBytes)
	resOffset := extismKVList(mem.Offset())
	resMem := pdk.FindMemory(resOffset)
	resBytes := resMem.ReadBytes()

	var res PluginKVListResponse
	if err := json.Unmarshal(resBytes, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !res.Result {
		return nil, fmt.Errorf("kv_list error: %s", res.Error)
	}

	return res.Keys, nil
}
