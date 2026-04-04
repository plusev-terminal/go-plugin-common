package tax

// Trade action constants (PluginRecord.Action for record_type "trade").
const (
	ActionBuy  = 0
	ActionSell = 1
)

// Transfer action constants (PluginRecord.Action for record_type "transfer").
const (
	ActionDeposit    = 0
	ActionWithdrawal = 1
)

// PluginGetRecordsRequest is the input the plugin passes to get_records.
type PluginGetRecordsRequest struct {
	From     string `json:"from"`      // RFC3339
	To       string `json:"to"`        // RFC3339
	Page     int    `json:"page"`      // 1-based
	PageSize int    `json:"page_size"` // max 1000; 0 -> 200
}

// PluginGetRecordsResponse is returned to the plugin from get_records.
type PluginGetRecordsResponse struct {
	Result     bool           `json:"result"`
	Records    []PluginRecord `json:"records,omitempty"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
	Error      string         `json:"error,omitempty"`
}

// PluginRecord is a flattened view of a HistoryRecord for plugins.
// All numeric fields are strings to preserve precision across the WASM boundary.
type PluginRecord struct {
	ID               uint64 `json:"id"`
	UID              string `json:"uid"` // Globally unique: "trade:<id>" or "transfer:<id>"
	TxID             string `json:"tx_id"`
	Ts               string `json:"ts"` // RFC3339
	Account          string `json:"account"`
	Comment          string `json:"comment"`
	RecordType       string `json:"record_type"` // "trade" or "transfer"
	Ticker           string `json:"ticker"`      // Trade only
	Asset            string `json:"asset"`
	Quote            string `json:"quote"` // Trade only
	Amount           string `json:"amount"`
	Action           int    `json:"action"`        // Trade: ActionBuy(0)/ActionSell(1); Transfer: ActionDeposit(0)/ActionWithdrawal(1)
	OrderType        int    `json:"order_type"`    // Trade only (-1 for transfers)
	Price            string `json:"price"`         // Trade: asset price in quote currency; Transfer: empty
	PriceC           string `json:"price_c"`       // Asset price in cost-basis currency (e.g. EUR)
	Value            string `json:"value"`         // Trade: total value in quote; Transfer: empty
	ValueC           string `json:"value_c"`       // Total value in cost-basis currency
	QuotePriceC      string `json:"quote_price_c"` // Trade only: quote currency price in cost-basis currency
	Fee              string `json:"fee"`           // Fee amount in fee_currency
	FeeCurrency      string `json:"fee_currency"`
	FeeC             string `json:"fee_c"`              // Fee amount in cost-basis currency
	FeePriceC        string `json:"fee_price_c"`        // Fee currency price in cost-basis currency
	QuoteFeeAmount   string `json:"quote_fee_amount"`   // Trade only: fee in quote currency
	QuoteFeeCurrency string `json:"quote_fee_currency"` // Trade only
	Source           string `json:"source"`             // Transfer only: origin account
	Destination      string `json:"destination"`        // Transfer only: target account
	TransferCategory string `json:"transfer_category"`  // Transfer only: mining, staking_reward, airdrop, etc.
	IsMarginTrade    bool   `json:"is_margin_trade"`    // Trade only
	IsDerivative     bool   `json:"is_derivative"`      // Trade only
	IsPhysical       bool   `json:"is_physical"`        // Trade only
	Status           string `json:"status"`
}

// PluginSettingsResponse is returned to the plugin from get_settings.
type PluginSettingsResponse struct {
	Result            bool   `json:"result"`
	CostBasisCurrency string `json:"cost_basis_currency"`
	Timezone          string `json:"timezone"`
	DateTimeFormat    string `json:"date_time_format"`
	Error             string `json:"error,omitempty"`
}

// PluginReportEntry is the input the plugin passes to submit_report_entry.
type PluginReportEntry struct {
	TxID          string `json:"tx_id"`
	RecordType    string `json:"record_type"`
	Ts            string `json:"ts"` // RFC3339
	Asset         string `json:"asset"`
	Amount        string `json:"amount"`
	PnL           string `json:"pnl"`
	HoldingPeriod int    `json:"holding_period"` // Days
	TaxCategory   string `json:"tax_category"`
	TaxableAmount string `json:"taxable_amount"`
	CostBasis     string `json:"cost_basis"`
	Proceeds      string `json:"proceeds"`
	Details       string `json:"details"` // Free-form JSON string
}

// PluginReportProgress is the input the plugin passes to report_progress.
type PluginReportProgress struct {
	Processed int    `json:"processed"`
	Total     int    `json:"total"`
	Phase     string `json:"phase"`   // e.g. "processing", "calculating", "finalizing"
	Message   string `json:"message"` // Optional status message
}

// ReportSummaryRow is a single key-value row in the report summary.
// Plugins produce country-specific summary metrics using this generic format.
type ReportSummaryRow struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Order int    `json:"order"`
}

// PluginReportSummary is the input the plugin passes to submit_report_summary.
type PluginReportSummary struct {
	Rows []ReportSummaryRow `json:"rows"`
}

// PluginKVPutRequest is the input for kv_put.
type PluginKVPutRequest struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     []byte `json:"value"` // Base64-encoded in JSON
}

// PluginKVGetRequest is the input for kv_get.
type PluginKVGetRequest struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
}

// PluginKVGetResponse is the output for kv_get.
type PluginKVGetResponse struct {
	Result bool   `json:"result"`
	Value  []byte `json:"value,omitempty"` // Base64-encoded in JSON
	Found  bool   `json:"found"`
	Error  string `json:"error,omitempty"`
}

// PluginKVDeleteRequest is the input for kv_delete.
type PluginKVDeleteRequest struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
}

// PluginKVListRequest is the input for kv_list.
type PluginKVListRequest struct {
	Namespace string `json:"namespace"`
	Prefix    string `json:"prefix"`
}

// PluginKVListResponse is the output for kv_list.
type PluginKVListResponse struct {
	Result bool     `json:"result"`
	Keys   []string `json:"keys,omitempty"`
	Error  string   `json:"error,omitempty"`
}
