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
