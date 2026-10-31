package meta

import "time"

// RateLimitScope selects how the host shares a rule's bucket between
// callers. Mirrors the host's ratelimit.Scope; the JSON values are the wire
// contract.
type RateLimitScope string

const (
	// RateLimitScopeIP is one bucket per host shared by all callers
	// (default when Scope is empty).
	RateLimitScopeIP RateLimitScope = "ip"

	// RateLimitScopeCredential is one bucket per host + connection
	// credential.
	RateLimitScopeCredential RateLimitScope = "credential"
)

// RateLimitRule declares one outbound rate limit for the host of the
// NetworkTargetRule pattern it rides on. It is a mirror of the host's
// ratelimit.Rule wire format. Host semantics:
//
//   - Rules apply to the pattern's HOST, not to this exact pattern.
//   - A host with several declared rules only passes a request when every
//     rule passes.
//   - The same rule ID declared by several plugins collapses into one
//     rule; the strictest declaration wins.
//   - Users can override everything in Terminal (System → Rate Limits).
//   - A 429/Retry-After returned by the upstream host pauses the host for
//     every caller, not just this plugin.
//   - Wildcard-host patterns (*.x.com) cannot carry limits; the host
//     ignores limits declared on them.
type RateLimitRule struct {
	ID            string         `json:"id"`                 // unique within a host, e.g. "total", "orders"
	Label         string         `json:"label,omitempty"`    // UI label, e.g. "Request weight"
	Limit         int            `json:"limit"`              // units per window, > 0
	WindowMs      int64          `json:"windowMs"`           // > 0
	Burst         int            `json:"burst,omitempty"`    // bucket capacity; 0 → Limit
	Scope         RateLimitScope `json:"scope,omitempty"`    // "" → RateLimitScopeIP
	PerKeyEnabled bool           `json:"perKey,omitempty"`   // one bucket per request rate key (see requester/types.Request.RateKey)
	CostWeighted  bool           `json:"weighted,omitempty"` // consume request Cost units instead of 1 (see requester/types.Request.Cost)
	Paths         []string       `json:"paths,omitempty"`    // URL path prefixes the rule applies to; empty = all
}

// NewRateLimitRule creates a rule allowing limit units per window, e.g.
// NewRateLimitRule("total", 1200, time.Minute) = 1200 units per minute.
func NewRateLimitRule(id string, limit int, window time.Duration) RateLimitRule {
	return RateLimitRule{
		ID:       id,
		Limit:    limit,
		WindowMs: window.Milliseconds(),
	}
}

// WithLabel sets the UI label shown in Terminal's Rate Limits page.
func (r RateLimitRule) WithLabel(label string) RateLimitRule {
	r.Label = label
	return r
}

// WithBurst sets the bucket capacity; 0 means the bucket holds Limit units.
func (r RateLimitRule) WithBurst(burst int) RateLimitRule {
	r.Burst = burst
	return r
}

// PerCredential gives each connection credential its own bucket instead of
// sharing one bucket per host.
func (r RateLimitRule) PerCredential() RateLimitRule {
	r.Scope = RateLimitScopeCredential
	return r
}

// PerKey gives each request rate key its own bucket (see
// requester/types.Request.RateKey), e.g. one bucket per JSON-RPC method.
func (r RateLimitRule) PerKey() RateLimitRule {
	r.PerKeyEnabled = true
	return r
}

// Weighted makes requests consume their Cost units (see
// requester/types.Request.Cost) instead of 1, e.g. an exchange's request
// weight.
func (r RateLimitRule) Weighted() RateLimitRule {
	r.CostWeighted = true
	return r
}

// ForPaths narrows the rule to requests whose URL path starts with one of
// the prefixes; without it the rule applies to all paths of the host.
func (r RateLimitRule) ForPaths(prefixes ...string) RateLimitRule {
	r.Paths = append([]string(nil), prefixes...)
	return r
}
