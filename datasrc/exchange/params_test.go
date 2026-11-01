package exchange

import (
	"testing"

	acct "github.com/plusev-terminal/go-plugin-common/v2/datasrc/exchange/account"
	tt "github.com/plusev-terminal/go-plugin-common/v2/trading"
)

func validOrder() CreateOrderParams {
	return CreateOrderParams{
		Market:   tt.Market{Symbol: "BTCUSDT", AssetType: tt.AssetTypePerpetual},
		Side:     acct.OrderSideBuy,
		Type:     acct.OrderTypeLimit,
		Quantity: "1",
		Price:    "42000",
	}
}

func TestCreateOrdersParamsValidate_MinimalLimit(t *testing.T) {
	p := CreateOrdersParams{Orders: []CreateOrderParams{validOrder()}}
	if err := p.Validate(); err != nil {
		t.Fatalf("valid limit order rejected: %v", err)
	}
}

func TestCreateOrdersParamsValidate_TriggerAliases(t *testing.T) {
	base := validOrder()
	base.Type = acct.OrderTypeStopMarket

	// TriggerPrice alone satisfies stop types.
	stop := base
	stop.TriggerPrice = "40000"
	if err := (CreateOrdersParams{Orders: []CreateOrderParams{stop}}).Validate(); err != nil {
		t.Fatalf("stop with triggerPrice rejected: %v", err)
	}

	// The type-specific alias satisfies it too.
	aliased := base
	aliased.StopPrice = "40000"
	if err := (CreateOrdersParams{Orders: []CreateOrderParams{aliased}}).Validate(); err != nil {
		t.Fatalf("stop with stopPrice alias rejected: %v", err)
	}

	// Neither present → rejected.
	if err := (CreateOrdersParams{Orders: []CreateOrderParams{base}}).Validate(); err == nil {
		t.Fatal("stop without any trigger price must be rejected")
	}
}

func TestCreateOrdersParamsValidate_TrailingUnits(t *testing.T) {
	base := validOrder()
	base.Type = acct.OrderTypeTrailingStop

	amount := base
	amount.TrailingAmount = "500"
	if err := (CreateOrdersParams{Orders: []CreateOrderParams{amount}}).Validate(); err != nil {
		t.Fatalf("trailing with absolute amount rejected: %v", err)
	}

	percent := base
	percent.TrailingPercent = "1.5"
	if err := (CreateOrdersParams{Orders: []CreateOrderParams{percent}}).Validate(); err != nil {
		t.Fatalf("trailing with percent rejected: %v", err)
	}

	if err := (CreateOrdersParams{Orders: []CreateOrderParams{base}}).Validate(); err == nil {
		t.Fatal("trailing without amount or percent must be rejected")
	}
}

func TestCreateOrdersParamsValidate_GTDRequiresExpireTime(t *testing.T) {
	gtd := validOrder()
	gtd.TimeInForce = acct.TimeInForceGTD
	if err := (CreateOrdersParams{Orders: []CreateOrderParams{gtd}}).Validate(); err == nil {
		t.Fatal("GTD without expireTime must be rejected")
	}

	gtd.ExpireTime = "1735689600000"
	if err := (CreateOrdersParams{Orders: []CreateOrderParams{gtd}}).Validate(); err != nil {
		t.Fatalf("GTD with expireTime rejected: %v", err)
	}
}

func TestCreateOrdersParamsValidate_AssetTypeRequired(t *testing.T) {
	for _, at := range []tt.AssetType{"", "linear", "futures "} {
		o := validOrder()
		o.Market.AssetType = at
		if err := (CreateOrdersParams{Orders: []CreateOrderParams{o}}).Validate(); err == nil {
			t.Errorf("order with assetType %q must be rejected", at)
		}
	}
}

func TestCreateOrdersParamsValidate_OptionalFieldRules(t *testing.T) {
	both := validOrder()
	both.Type = acct.OrderTypeTrailingStop
	both.Price = ""
	both.TrailingAmount = "50"
	both.TrailingPercent = "1"

	expiryWithoutGTD := validOrder()
	expiryWithoutGTD.ExpireTime = "1893456000000"

	unknownSource := validOrder()
	unknownSource.TriggerSource = "MARK_PRICE"

	for name, o := range map[string]CreateOrderParams{
		"both trailing units":    both,
		"expireTime without gtd": expiryWithoutGTD,
		"unknown triggerSource":  unknownSource,
	} {
		if err := (CreateOrdersParams{Orders: []CreateOrderParams{o}}).Validate(); err == nil {
			t.Errorf("%s: must be rejected", name)
		}
	}
}

func TestCancelAndModifyRequireAssetTypeWithMarket(t *testing.T) {
	if err := (CancelOrdersParams{}).Validate(); err != nil {
		t.Fatalf("cancel-all without a market stays allowed: %v", err)
	}
	if err := (CancelOrdersParams{Market: tt.Market{Symbol: "BTCUSDT"}}).Validate(); err == nil {
		t.Fatal("cancel on a market without assetType must be rejected")
	}
	if err := (CancelOrdersParams{Market: tt.Market{Symbol: "BTCUSDT", AssetType: tt.AssetTypeSpot}, IDs: []string{"1"}}).Validate(); err != nil {
		t.Fatalf("typed cancel rejected: %v", err)
	}

	modify := ModifyOrdersParams{Market: tt.Market{Symbol: "BTCUSDT"}, Orders: []ModifyOrderParams{{OrderID: "1", Price: "1"}}}
	if err := modify.Validate(); err == nil {
		t.Fatal("modify without assetType must be rejected")
	}
	modify.Market.AssetType = tt.AssetTypePerpetual
	if err := modify.Validate(); err != nil {
		t.Fatalf("typed modify rejected: %v", err)
	}
}

func TestMarketDataRejectsUnknownAssetType(t *testing.T) {
	ohlcv := GetOHLCVParams{Market: tt.Market{Symbol: "BTCUSDT"}, Timeframe: "1m"}
	if err := ohlcv.Validate(); err != nil {
		t.Fatalf("unspecified assetType is left to the venue: %v", err)
	}
	ohlcv.Market.AssetType = "linear"
	if err := ohlcv.Validate(); err == nil {
		t.Fatal("unknown assetType must be rejected")
	}

	stream := OHLCVStreamParams{Market: tt.Market{Symbol: "BTCUSDT", AssetType: "perp"}, Timeframe: "1m"}
	if err := stream.Validate(); err == nil {
		t.Fatal("unknown stream assetType must be rejected")
	}
}

func TestParseCreateOrdersParamsStrict(t *testing.T) {
	order := map[string]any{
		"market":   map[string]any{"symbol": "BTCUSDT", "assetType": "perpetual", "category": "Perpetual"},
		"side":     "buy",
		"type":     "limit",
		"quantity": "1",
		"price":    "42000",
	}

	p, err := ParseCreateOrdersParams(map[string]any{"orders": []any{order}, "cancelGroupOnFailure": true})
	if err != nil || len(p.Orders) != 1 || !p.CancelGroupOnFailure || p.Orders[0].Market.AssetType != tt.AssetTypePerpetual {
		t.Fatalf("valid batch: got %+v, %v", p, err)
	}

	retired := map[string]any{"market": order["market"], "side": "sell", "type": "trailing_stop", "quantity": "1", "trailingDelta": "1"}
	if _, err := ParseCreateOrdersParams(map[string]any{"orders": []any{order, retired}}); err == nil {
		t.Fatal("a retired v1 field (trailingDelta) must be rejected, not dropped")
	}

	flat := map[string]any{"symbol": "BTCUSDT", "side": "buy", "type": "market", "quantity": "1"}
	if _, err := ParseCreateOrdersParams(map[string]any{"orders": []any{flat}}); err == nil {
		t.Fatal("a flat symbol key must be rejected (the contract carries market.symbol)")
	}

	if _, err := ParseCreateOrdersParams(map[string]any{"orders": []any{order, "not-an-object"}}); err == nil {
		t.Fatal("an unparsable order must fail the batch, not vanish from it")
	}

	if _, err := ParseCreateOrdersParams(map[string]any{"orders": []any{order}, "scope": "futures"}); err == nil {
		t.Fatal("unknown batch-level keys must be rejected")
	}
}
