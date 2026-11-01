package trading

import "testing"

func TestMarketKey(t *testing.T) {
	spot := Market{AssetType: AssetTypeSpot, Symbol: "BTCUSDT"}
	perp := Market{AssetType: AssetTypePerpetual, Symbol: "BTCUSDT"}

	if spot.Key() == perp.Key() {
		t.Fatalf("spot and perpetual markets with the same symbol must have distinct keys: %q", spot.Key())
	}
	if got, want := perp.Key(), "perpetual:BTCUSDT"; got != want {
		t.Fatalf("perp key = %q, want %q", got, want)
	}
	if got, want := spot.Key(), "spot:BTCUSDT"; got != want {
		t.Fatalf("spot key = %q, want %q", got, want)
	}
}

func TestAssetTypeIsDerivative(t *testing.T) {
	derivatives := []AssetType{AssetTypePerpetual, AssetTypeFutures, AssetTypeInversePerpetual, AssetTypeInverseFutures}
	for _, at := range derivatives {
		if !at.IsDerivative() {
			t.Errorf("%s must classify as derivative", at)
		}
	}
	if AssetTypeSpot.IsDerivative() || AssetTypeOption.IsDerivative() {
		t.Error("spot and option must not classify as derivative")
	}
}

func TestMarketKeyIncludesAddress(t *testing.T) {
	a := Market{AssetType: AssetTypeSpot, Symbol: "SOL/USDC", Address: "PoolA"}
	b := Market{AssetType: AssetTypeSpot, Symbol: "SOL/USDC", Address: "PoolB"}

	if a.Key() == b.Key() {
		t.Fatalf("two pools sharing a symbol must have distinct keys: %q", a.Key())
	}
	if got, want := a.Key(), "spot:SOL/USDC@PoolA"; got != want {
		t.Fatalf("key = %q, want %q", got, want)
	}
}

func TestAssetTypeValid(t *testing.T) {
	for _, at := range AssetTypes() {
		if !at.Valid() {
			t.Errorf("%s must be valid", at)
		}
	}
	for _, at := range []AssetType{"", "linear", "futures ", "Perpetual", "inverse"} {
		if at.Valid() {
			t.Errorf("%q must not be valid", at)
		}
	}
}

func TestResolveAssetType(t *testing.T) {
	cases := []struct {
		name      string
		requested AssetType
		listed    []AssetType
		want      AssetType
		wantErr   bool
	}{
		{"empty single-type venue resolves", "", []AssetType{AssetTypePerpetual}, AssetTypePerpetual, false},
		{"empty multi-type venue is ambiguous", "", []AssetType{AssetTypeSpot, AssetTypePerpetual}, "", true},
		{"empty without listing", "", nil, "", true},
		{"listed type passes", AssetTypeSpot, []AssetType{AssetTypeSpot, AssetTypePerpetual}, AssetTypeSpot, false},
		{"unlisted type rejected", AssetTypeOption, []AssetType{AssetTypeSpot, AssetTypePerpetual}, "", true},
		{"unknown type rejected", "linear", []AssetType{AssetTypePerpetual}, "", true},
		{"valid type without listing", AssetTypeFutures, nil, AssetTypeFutures, false},
	}

	for _, c := range cases {
		got, err := ResolveAssetType(c.requested, c.listed...)
		if (err != nil) != c.wantErr || got != c.want {
			t.Errorf("%s: got (%q, %v), want (%q, err=%v)", c.name, got, err, c.want, c.wantErr)
		}
	}
}

func TestFindMarket(t *testing.T) {
	markets := []Market{
		{Symbol: "BTCUSDT", AssetType: AssetTypeSpot},
		{Symbol: "BTCUSDT", AssetType: AssetTypePerpetual},
		{Symbol: "ETHUSDT", AssetType: AssetTypePerpetual},
		{Symbol: "SOL/USDC", AssetType: AssetTypeSpot, Address: "PoolA"},
		{Symbol: "SOL/USDC", AssetType: AssetTypeSpot, Address: "PoolB"},
	}

	if m, err := FindMarket(markets, Market{Symbol: "BTCUSDT", AssetType: AssetTypePerpetual}); err != nil || m.AssetType != AssetTypePerpetual {
		t.Fatalf("typed reference: got %+v, %v", m, err)
	}
	if _, err := FindMarket(markets, Market{Symbol: "BTCUSDT"}); err == nil {
		t.Fatal("symbol-only reference to a spot+perp symbol must be ambiguous")
	}
	if m, err := FindMarket(markets, Market{Symbol: "ETHUSDT"}); err != nil || m.AssetType != AssetTypePerpetual {
		t.Fatalf("unique symbol resolves without assetType: got %+v, %v", m, err)
	}
	if m, err := FindMarket(markets, Market{Symbol: "SOL/USDC", Address: "PoolB"}); err != nil || m.Address != "PoolB" {
		t.Fatalf("address narrows pools: got %+v, %v", m, err)
	}
	if _, err := FindMarket(markets, Market{Symbol: "SOL/USDC"}); err == nil {
		t.Fatal("two pools under one symbol must be ambiguous")
	}
	if _, err := FindMarket(markets, Market{Symbol: "BTCUSDT", AssetType: AssetTypeFutures}); err == nil {
		t.Fatal("missing market must error")
	}
	if _, err := FindMarket(markets, Market{}); err == nil {
		t.Fatal("empty reference must error")
	}
}

func TestUsesContracts(t *testing.T) {
	if (Market{}).UsesContracts() || (Market{QuantityDenomination: QuantityBase}).UsesContracts() {
		t.Fatal("base denomination is the default")
	}
	if !(Market{QuantityDenomination: QuantityContracts}).UsesContracts() {
		t.Fatal("contracts denomination must report UsesContracts")
	}
}
