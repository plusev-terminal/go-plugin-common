package tax

// PluginTrade is the standardized trade record that tax plugins submit to the host
// via the submit_trade host function. All numeric fields are strings to preserve precision.
type PluginTrade struct {
	TxID                  string `json:"tx_id"`
	Ts                    string `json:"ts"` // RFC3339
	Account               string `json:"account"`
	Comment               string `json:"comment"`
	Ticker                string `json:"ticker"`
	Quote                 string `json:"quote"`
	QuoteAddress          string `json:"quote_address"`
	Asset                 string `json:"asset"`
	AssetAddress          string `json:"asset_address"`
	Price                 string `json:"price"`
	PriceC                string `json:"price_c"`
	PriceConvertedBy      string `json:"price_converted_by"`
	QuotePriceC           string `json:"quote_price_c"`
	QuotePriceConvertedBy string `json:"quote_price_converted_by"`
	Amount                string `json:"amount"`
	Value                 string `json:"value"`
	ValueC                string `json:"value_c"`
	Action                string `json:"action"`     // BUY or SELL
	OrderType             string `json:"order_type"` // TAKER or MAKER
	OrderID               string `json:"order_id"`
	FeeAmount             string `json:"fee_amount"`
	FeeCurrency           string `json:"fee_currency"`
	FeeCurrencyAddress    string `json:"fee_currency_address"`
	FeeAmountC            string `json:"fee_amount_c"`
	FeePriceC             string `json:"fee_price_c"`
	FeeConvertedBy        string `json:"fee_converted_by"`
	FeeDecimals           int32  `json:"fee_decimals"`
	AssetDecimals         int32  `json:"asset_decimals"`
	QuoteDecimals         int32  `json:"quote_decimals"`
	IsMarginTrade         bool   `json:"is_margin_trade"`
	IsDerivative          bool   `json:"is_derivative"`
	IsPhysical            bool   `json:"is_physical"`
}

// PluginTransfer is the standardized transfer record that tax plugins submit to the host
// via the submit_transfer host function.
type PluginTransfer struct {
	TxID               string `json:"tx_id"`
	Ts                 string `json:"ts"`
	Account            string `json:"account"`
	Source             string `json:"source"`
	Destination        string `json:"destination"`
	Comment            string `json:"comment"`
	Asset              string `json:"asset"`
	AssetAddress       string `json:"asset_address"`
	AssetDecimals      int32  `json:"asset_decimals"`
	Amount             string `json:"amount"`
	Action             string `json:"action"` // DEPOSIT or WITHDRAWAL
	Fee                string `json:"fee"`
	FeeDecimals        int32  `json:"fee_decimals"`
	FeeCurrency        string `json:"fee_currency"`
	FeeCurrencyAddress string `json:"fee_currency_address"`
	FeeC               string `json:"fee_c"`
	FeeConvertedBy     string `json:"fee_converted_by"`
	FeePriceC          string `json:"fee_price_c"`
	TransferCategory   string `json:"transfer_category"`
}

// PluginSubmitResult is returned by the host after a submit_trade or submit_transfer call.
type PluginSubmitResult struct {
	Result bool   `json:"result"`
	ID     uint64 `json:"id,omitempty"`
	Error  string `json:"error,omitempty"`
}
