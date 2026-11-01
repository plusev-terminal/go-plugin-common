package exchange

import (
	"encoding/json"
	"testing"

	acct "github.com/plusev-terminal/go-plugin-common/v2/datasrc/exchange/account"
	tt "github.com/plusev-terminal/go-plugin-common/v2/trading"
)

func declared() *OrderCapabilities {
	return &OrderCapabilities{
		AssetTypes:           []tt.AssetType{tt.AssetTypePerpetual},
		OrderTypes:           []acct.OrderType{acct.OrderTypeLimit, acct.OrderTypeMarket, acct.OrderTypeStopMarket},
		TimeInForce:          []acct.TimeInForce{acct.TimeInForceGTC, acct.TimeInForceIOC},
		TriggerSources:       []acct.TriggerSource{acct.TriggerSourceMark},
		DefaultTriggerSource: acct.TriggerSourceMark,
		ReduceOnly:           true,
		ClientOrderID:        true,
		ExtraKeys:            []string{"venueFlag"},
	}
}

func TestCapabilitiesAcceptDeclaredFeatures(t *testing.T) {
	o := validOrder()
	o.TimeInForce = acct.TimeInForceIOC
	o.ReduceOnly = true
	o.ClientOrderID = "glf_1_entry"
	o.Extra = map[string]any{"venueFlag": true}

	if err := declared().Check(o); err != nil {
		t.Fatalf("declared features rejected: %v", err)
	}
}

func TestCapabilitiesRejectUndeclaredFeatures(t *testing.T) {
	mutations := map[string]func(*CreateOrderParams){
		"asset type":     func(o *CreateOrderParams) { o.Market.AssetType = tt.AssetTypeSpot },
		"order type":     func(o *CreateOrderParams) { o.Type = acct.OrderTypeTrailingStop },
		"time in force":  func(o *CreateOrderParams) { o.TimeInForce = acct.TimeInForceGTD },
		"trigger source": func(o *CreateOrderParams) { o.TriggerSource = acct.TriggerSourceIndex },
		"post only":      func(o *CreateOrderParams) { o.PostOnly = true },
		"trailing pct":   func(o *CreateOrderParams) { o.TrailingPercent = "1" },
		"position side":  func(o *CreateOrderParams) { o.PositionSide = "long" },
		"extra key":      func(o *CreateOrderParams) { o.Extra = map[string]any{"triggerPriceType": "LAST_PRICE"} },
	}

	for name, mutate := range mutations {
		o := validOrder()
		mutate(&o)

		if err := declared().Check(o); err == nil {
			t.Errorf("%s: undeclared feature must be rejected", name)
		}
	}
}

func TestCapabilitiesDefaultTriggerSourceAlwaysPasses(t *testing.T) {
	caps := declared()
	caps.TriggerSources = nil

	o := validOrder()
	o.TriggerSource = acct.TriggerSourceMark

	if err := caps.Check(o); err != nil {
		t.Fatalf("asking explicitly for the venue default must pass: %v", err)
	}
}

func TestCapabilitiesZeroValueSupportsNothing(t *testing.T) {
	if err := (&OrderCapabilities{}).Check(validOrder()); err == nil {
		t.Fatal("a zero declaration must reject orders")
	}

	var nilCaps *OrderCapabilities
	if err := nilCaps.Check(validOrder()); err == nil {
		t.Fatal("no declaration must reject orders")
	}
}

func TestLegacyBaseline(t *testing.T) {
	legacy := LegacyOrderCapabilities()

	o := validOrder()
	o.Market.AssetType = tt.AssetTypeSpot
	o.ReduceOnly = true
	o.Extra = map[string]any{"anything": 1}

	if err := legacy.Check(o); err != nil {
		t.Fatalf("the v1 surface must stay allowed for pre-v2.1.0 plugins: %v", err)
	}

	for name, mutate := range map[string]func(*CreateOrderParams){
		"trigger source":  func(o *CreateOrderParams) { o.TriggerSource = acct.TriggerSourceMark },
		"gtd":             func(o *CreateOrderParams) { o.TimeInForce = acct.TimeInForceGTD; o.ExpireTime = "1" },
		"trailing amount": func(o *CreateOrderParams) { o.Type = acct.OrderTypeTrailingStop; o.TrailingAmount = "5" },
	} {
		o := validOrder()
		mutate(&o)

		if err := legacy.Check(o); err == nil {
			t.Errorf("%s: a pre-v2.1.0 binary cannot read it — must be rejected", name)
		}
	}

	if err := legacy.CheckAll(CreateOrdersParams{Orders: []CreateOrderParams{validOrder()}, CancelGroupOnFailure: true}); err == nil {
		t.Fatal("group cancellation was never implemented by v1 plugins — must be rejected")
	}
}

func TestCapabilitiesJSONRoundTripDropsLegacyFlag(t *testing.T) {
	raw, err := json.Marshal(LegacyOrderCapabilities())
	if err != nil {
		t.Fatal(err)
	}

	var decoded OrderCapabilities
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Legacy() {
		t.Fatal("the legacy marker must never travel over the wire")
	}
}
