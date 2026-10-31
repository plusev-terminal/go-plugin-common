package meta

import "sort"

type Meta struct {
	PluginID    string          `json:"pluginId" validate:"required"`
	Name        string          `json:"name" validate:"required"`
	AppID       string          `json:"appId" validate:"required"`
	Category    string          `json:"category"` // Plugin category specific to the app id.
	Description string          `json:"description"`
	Author      string          `json:"author"`
	Version     string          `json:"version" validate:"required,semver"`
	Repository  string          `json:"repository"`
	Tags        []string        `json:"tags"`
	Contacts    []AuthorContact `json:"contacts"`
	Resources   ResourceAccess  `json:"resources"`
	Features    []string        `json:"features"` // List of supported features
	Chain       string          `json:"chain"`    // Blockchain identifier for chain-category plugins
	// Commands optionally declares a rate-limit priority class per command;
	// see CommandInfo. Nil/empty is omitted, so plugins that declare no
	// priorities produce byte-identical meta JSON to v1.15.0.
	Commands []CommandInfo `json:"commands,omitempty"`
}

// CommandInfo declares one plugin command and its rate-limit priority class
// for the host's dispatch scheduling; see the Priority* constants.
type CommandInfo struct {
	Name string `json:"name"`
	// Priority should be one of the Priority* constants. NOT validated here:
	// unknown values pass through the wire as-is and are ignored by the host
	// (which falls back to its canonical command-name table), never fatal.
	Priority string `json:"priority,omitempty"`
}

// Rate-limit priority classes a plugin may declare for its commands.
const (
	PriorityTrading     = "trading"
	PriorityInteractive = "interactive"
	PriorityBackground  = "background"
)

// SortCommandInfosByName returns a copy of infos ordered by Name so the meta
// export is deterministic regardless of map iteration order. It lives in the
// natively testable meta package because the plugin package (its caller) is
// wasm-only.
func SortCommandInfosByName(infos []CommandInfo) []CommandInfo {
	sorted := make([]CommandInfo, len(infos))
	copy(sorted, infos)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
	return sorted
}

type AuthorContact struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type ResourceAccess struct {
	AllowedNetworkTargets []NetworkTargetRule `json:"allowedNetworkTargets"`
	FsWriteAccess         map[string]string   `json:"fsWriteAccess"`
}

type NetworkTargetRule struct {
	Pattern string `json:"pattern"`
	// Limits optionally declares outbound rate limits enforced by the host
	// for this pattern's HOST; see RateLimitRule for the semantics. The
	// host ignores limits on wildcard-host patterns (*.x.com).
	Limits         []RateLimitRule `json:"limits,omitempty"`
	TradingReserve *float64        `json:"tradingReserve,omitempty"`
}
