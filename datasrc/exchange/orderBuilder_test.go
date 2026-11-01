package exchange

import (
	"testing"

	acct "github.com/plusev-terminal/go-plugin-common/v2/datasrc/exchange/account"
	tt "github.com/plusev-terminal/go-plugin-common/v2/trading"
)

func TestBuildGroupAppliesTriggerSourceToProtection(t *testing.T) {
	market := tt.Market{Symbol: "BTCUSDT", AssetType: tt.AssetTypePerpetual}

	orders, err := NewOrder(market, acct.OrderSideBuy).
		Market("1").
		TriggerSource(acct.TriggerSourceMark).
		WithStopLoss("40000").
		WithTakeProfit("50000").
		BuildGroup()
	if err != nil {
		t.Fatal(err)
	}

	if orders[0].TriggerSource != "" {
		t.Fatalf("a market entry has nothing to trigger, got %q", orders[0].TriggerSource)
	}
	for _, leg := range orders[1:] {
		if leg.TriggerSource != acct.TriggerSourceMark {
			t.Fatalf("%s leg must trigger on mark, got %q", leg.Type, leg.TriggerSource)
		}
	}
	if err := (CreateOrdersParams{Orders: orders}).Validate(); err != nil {
		t.Fatalf("the bracket must validate: %v", err)
	}
}

func TestBuilderNewFields(t *testing.T) {
	market := tt.Market{Symbol: "BTCUSDT", AssetType: tt.AssetTypePerpetual}

	trailing := NewOrder(market, acct.OrderSideSell).TrailingStopPercent("1", "1.5").Build()
	if trailing.Type != acct.OrderTypeTrailingStop || trailing.TrailingPercent != "1.5" || trailing.TrailingAmount != "" {
		t.Fatalf("percent trailing stop: %+v", trailing)
	}

	gtd := NewOrder(market, acct.OrderSideBuy).Limit("1", "42000").ExpireAt(1893456000000).Build()
	if gtd.TimeInForce != acct.TimeInForceGTD || gtd.ExpireTime != "1893456000000" {
		t.Fatalf("good-till-date: %+v", gtd)
	}

	for _, o := range []CreateOrderParams{trailing, gtd} {
		if err := (CreateOrdersParams{Orders: []CreateOrderParams{o}}).Validate(); err != nil {
			t.Fatalf("%s must validate: %v", o.Type, err)
		}
	}
}
