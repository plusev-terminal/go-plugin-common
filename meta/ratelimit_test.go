package meta

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNetworkTargetRule_NoLimitsJSONUnchanged(t *testing.T) {
	data, err := json.Marshal(NetworkTargetRule{Pattern: "https://x/*"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Must stay byte-identical to v1.14.0: the new fields are omitempty.
	want := `{"pattern":"https://x/*"}`
	if string(data) != want {
		t.Fatalf("Expected %s, got %s", want, string(data))
	}
}

func TestNetworkTargetRule_FullGoldenJSON(t *testing.T) {
	reserve := 0.25
	rule := NetworkTargetRule{
		Pattern: "https://api.example.com/v1/*",
		Limits: []RateLimitRule{
			NewRateLimitRule("total", 1200, time.Minute).
				WithLabel("Request weight").
				WithBurst(2400).
				PerCredential().
				PerKey().
				Weighted().
				ForPaths("/api/v5/public", "/api/v5/trade"),
		},
		TradingReserve: &reserve,
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	want := `{"pattern":"https://api.example.com/v1/*","limits":[{"id":"total","label":"Request weight","limit":1200,"windowMs":60000,"burst":2400,"scope":"credential","perKey":true,"weighted":true,"paths":["/api/v5/public","/api/v5/trade"]}],"tradingReserve":0.25}`
	if string(data) != want {
		t.Fatalf("Expected %s, got %s", want, string(data))
	}
}

func TestNetworkTargetRule_TradingReserveOmittedWhenNil(t *testing.T) {
	rule := NetworkTargetRule{
		Pattern: "https://api.example.com/v1/*",
		Limits: []RateLimitRule{
			NewRateLimitRule("orders", 10, time.Second),
		},
	}

	data, err := json.Marshal(rule)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	want := `{"pattern":"https://api.example.com/v1/*","limits":[{"id":"orders","limit":10,"windowMs":1000}]}`
	if string(data) != want {
		t.Fatalf("Expected %s, got %s", want, string(data))
	}
}

func TestRateLimitRuleHelpers(t *testing.T) {
	base := NewRateLimitRule("orders", 10, 5*time.Second)
	if base.ID != "orders" || base.Limit != 10 || base.WindowMs != 5000 {
		t.Fatalf("Unexpected base rule: %+v", base)
	}
	if base.Scope != "" || base.Burst != 0 || base.PerKeyEnabled || base.CostWeighted || base.Paths != nil {
		t.Fatalf("Expected zero-valued optionals, got %+v", base)
	}

	labeled := base.WithLabel("Orders")
	if labeled.Label != "Orders" {
		t.Fatalf("Expected label Orders, got %q", labeled.Label)
	}
	if base.Label != "" {
		t.Fatalf("Helpers must return copies; base was mutated: %+v", base)
	}

	burst := labeled.WithBurst(20)
	if burst.Burst != 20 {
		t.Fatalf("Expected burst 20, got %d", burst.Burst)
	}

	keyed := burst.PerKey()
	if !keyed.PerKeyEnabled {
		t.Fatal("Expected PerKey to set PerKeyEnabled")
	}

	weighted := keyed.Weighted()
	if !weighted.CostWeighted {
		t.Fatal("Expected Weighted to set CostWeighted")
	}

	cred := weighted.PerCredential()
	if cred.Scope != RateLimitScopeCredential {
		t.Fatalf("Expected scope credential, got %q", cred.Scope)
	}

	paths := cred.ForPaths("/a", "/b")
	if len(paths.Paths) != 2 || paths.Paths[0] != "/a" || paths.Paths[1] != "/b" {
		t.Fatalf("Unexpected paths: %v", paths.Paths)
	}
	if base.Paths != nil || cred.Paths != nil {
		t.Fatal("ForPaths must not leak into earlier copies")
	}
}
