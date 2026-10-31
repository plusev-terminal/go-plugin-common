package meta

import (
	"encoding/json"
	"testing"
)

// The golden string below was produced by marshaling this exact Meta against
// the untouched v1.15.0 struct definitions. Plugins that declare no command
// priorities must stay byte-identical to v1.15.0.
func goldenMeta() Meta {
	return Meta{
		PluginID: "com.example.plugin",
		Name:     "Example Plugin",
		AppID:    "datasrc",
		Version:  "1.0.0",
		Features: []string{"getMarkets", "createOrders"},
	}
}

const goldenMetaNoCommands = `{"pluginId":"com.example.plugin","name":"Example Plugin","appId":"datasrc","category":"","description":"","author":"","version":"1.0.0","repository":"","tags":null,"contacts":null,"resources":{"allowedNetworkTargets":null,"fsWriteAccess":null},"features":["getMarkets","createOrders"],"chain":""}`

func TestMeta_NoCommandsJSONUnchanged(t *testing.T) {
	data, err := json.Marshal(goldenMeta())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Must stay byte-identical to v1.15.0: Commands is omitempty.
	if string(data) != goldenMetaNoCommands {
		t.Fatalf("Expected %s, got %s", goldenMetaNoCommands, string(data))
	}
}

func TestMeta_WithCommandsGoldenJSON(t *testing.T) {
	m := goldenMeta()
	m.Commands = []CommandInfo{
		{Name: "createAlgoOrders", Priority: PriorityTrading},
		{Name: "fetchHistory", Priority: PriorityBackground},
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	want := `{"pluginId":"com.example.plugin","name":"Example Plugin","appId":"datasrc","category":"","description":"","author":"","version":"1.0.0","repository":"","tags":null,"contacts":null,"resources":{"allowedNetworkTargets":null,"fsWriteAccess":null},"features":["getMarkets","createOrders"],"chain":"","commands":[{"name":"createAlgoOrders","priority":"trading"},{"name":"fetchHistory","priority":"background"}]}`
	if string(data) != want {
		t.Fatalf("Expected %s, got %s", want, string(data))
	}
}

// Declared priorities travel the wire verbatim; the host validates them and
// falls back to its canonical name table, never fails.
func TestCommandInfo_InvalidPriorityPassesThrough(t *testing.T) {
	data, err := json.Marshal(CommandInfo{Name: "customCmd", Priority: "turbo"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	want := `{"name":"customCmd","priority":"turbo"}`
	if string(data) != want {
		t.Fatalf("Expected %s, got %s", want, string(data))
	}
}

func TestSortCommandInfosByName(t *testing.T) {
	in := []CommandInfo{
		{Name: "zeta", Priority: PriorityTrading},
		{Name: "alpha", Priority: PriorityBackground},
		{Name: "mid", Priority: PriorityInteractive},
	}
	got := SortCommandInfosByName(in)

	want := []string{"alpha", "mid", "zeta"}
	for i, name := range want {
		if got[i].Name != name {
			t.Fatalf("Expected %v, got %v", want, got)
		}
	}

	// Input must not be mutated.
	if in[0].Name != "zeta" {
		t.Fatalf("Expected input to be unmodified, got %v", in)
	}

	if got := SortCommandInfosByName(nil); len(got) != 0 {
		t.Fatalf("Expected empty result for nil input, got %v", got)
	}
}
