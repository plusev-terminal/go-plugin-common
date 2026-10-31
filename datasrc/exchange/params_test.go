package exchange

import (
	"testing"

	acct "github.com/plusev-terminal/go-plugin-common/v2/datasrc/exchange/account"
	tt "github.com/plusev-terminal/go-plugin-common/v2/trading"
)

func validOrder() CreateOrderParams {
	return CreateOrderParams{
		Market: tt.Market{Symbol: "BTCUSDT", AssetType: tt.AssetTypePerpetual},
		Side:   acct.OrderSideBuy,
		Type:   acct.OrderTypeLimit,
		Quantity: "1",
		Price:   "42000",
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
