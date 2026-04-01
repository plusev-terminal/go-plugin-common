package trading

// TradeRecord represents a single normalized trade (fill) from an exchange.
// All numeric fields are strings to preserve precision — consumers should use
// high-precision libraries like shopspring/decimal to parse these values.
//
// This is the common response type for getTradeHistory / getOrderHistory commands.
// Each datasrc plugin normalizes its exchange-specific trade format into this struct.
type TradeRecord struct {
	TxID        string `json:"txId"`                  // Exchange-assigned trade/fill ID (used for dedup)
	OrderID     string `json:"orderId,omitempty"`     // Exchange-assigned order ID
	Timestamp   int64  `json:"timestamp"`             // Unix milliseconds
	Symbol      string `json:"symbol"`                // Trading pair symbol, e.g. "BTC/USDT"
	Base        string `json:"base"`                  // Base asset, e.g. "BTC"
	Quote       string `json:"quote"`                 // Quote asset, e.g. "USDT"
	Side        string `json:"side"`                  // "buy" or "sell"
	Price       string `json:"price"`                 // Execution price in quote currency
	Amount      string `json:"amount"`                // Quantity of base asset
	Fee         string `json:"fee,omitempty"`         // Fee amount
	FeeCurrency string `json:"feeCurrency,omitempty"` // Currency of the fee
	OrderType   string `json:"orderType,omitempty"`   // "maker" or "taker"
	AssetType   string `json:"assetType,omitempty"`   // "spot", "perpetual", "futures", "option"
	Account     string `json:"account,omitempty"`     // Account or sub-account label

	// Precision hints (optional, may be 0 if not provided by the exchange)
	PricePrecision    int `json:"pricePrecision,omitempty"`
	QuantityPrecision int `json:"quantityPrecision,omitempty"`
	FeePrecision      int `json:"feePrecision,omitempty"`

	// Margin / derivatives metadata
	IsMarginTrade bool   `json:"isMarginTrade,omitempty"`
	IsDerivative  bool   `json:"isDerivative,omitempty"`
	IsPhysical    bool   `json:"isPhysical,omitempty"` // Physical delivery vs cash-settled
	ContractSize  string `json:"contractSize,omitempty"`
}

// TransferRecord represents a single normalized deposit or withdrawal from an exchange.
// All numeric fields are strings to preserve precision.
//
// This is the common response type for getTransferHistory commands.
type TransferRecord struct {
	TxID        string `json:"txId"`                  // Exchange-assigned transfer ID (used for dedup)
	Timestamp   int64  `json:"timestamp"`             // Unix milliseconds
	Asset       string `json:"asset"`                 // The asset being transferred, e.g. "BTC"
	Amount      string `json:"amount"`                // Transfer amount
	Action      string `json:"action"`                // "deposit" or "withdrawal"
	Fee         string `json:"fee,omitempty"`         // Fee amount
	FeeCurrency string `json:"feeCurrency,omitempty"` // Currency of the fee
	Source      string `json:"source,omitempty"`      // Source address/account
	Destination string `json:"destination,omitempty"` // Destination address/account
	Network     string `json:"network,omitempty"`     // Blockchain network (e.g. "ethereum", "solana")
	Status      string `json:"status,omitempty"`      // Exchange-reported status (e.g. "completed", "pending")
	Account     string `json:"account,omitempty"`     // Account or sub-account label

	// Precision hints
	AssetPrecision int `json:"assetPrecision,omitempty"`
	FeePrecision   int `json:"feePrecision,omitempty"`
}

// ---------------------------------------------------------------------------
// Paginated response wrappers
// ---------------------------------------------------------------------------
// Plugins return these from getTradeHistory / getOrderHistory / getTransferHistory.
// The host passes NextCursor back into the next request's Cursor param to fetch
// the next page. When HasMore is false the host stops paginating.

// PaginatedTradeResponse wraps a page of TradeRecords with cursor metadata.
type PaginatedTradeResponse struct {
	Records    []TradeRecord `json:"records"`
	NextCursor string        `json:"nextCursor,omitempty"` // Opaque; host passes back as Cursor param
	HasMore    bool          `json:"hasMore"`
}

// PaginatedTransferResponse wraps a page of TransferRecords with cursor metadata.
type PaginatedTransferResponse struct {
	Records    []TransferRecord `json:"records"`
	NextCursor string           `json:"nextCursor,omitempty"`
	HasMore    bool             `json:"hasMore"`
}
