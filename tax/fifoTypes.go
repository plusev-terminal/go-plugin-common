package tax

// FifoAddLotRequest is the input for fifo_add_lot.
type FifoAddLotRequest struct {
	Account   string         `json:"account"`
	Asset     string         `json:"asset"`
	ID        string         `json:"id"` // Caller-assigned identifier (e.g. TxID)
	Ts        string         `json:"ts"` // RFC3339
	Amount    string         `json:"amount"`
	CostBasis string         `json:"cost_basis"` // Per-unit, in cost-basis currency
	Fee       string         `json:"fee"`        // Per-unit, in cost-basis currency
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// FifoDisposeRequest is the input for fifo_dispose.
type FifoDisposeRequest struct {
	Account string `json:"account"`
	Asset   string `json:"asset"`
	Ts      string `json:"ts"` // RFC3339 — timestamp of the disposal event
	Amount  string `json:"amount"`
}

// FifoDisposeResponse is returned by the host from fifo_dispose.
type FifoDisposeResponse struct {
	Result          bool           `json:"result"`
	Matches         []FifoLotMatch `json:"matches,omitempty"`
	MatchedAmount   string         `json:"matched_amount"`
	UnmatchedAmount string         `json:"unmatched_amount"`
	Error           string         `json:"error,omitempty"`
}

// FifoLotMatch represents one lot (or partial lot) consumed by a disposal.
type FifoLotMatch struct {
	LotID       string `json:"lot_id"`
	LotTs       string `json:"lot_ts"` // RFC3339
	Amount      string `json:"amount"`
	CostBasis   string `json:"cost_basis"` // Total for matched amount
	Fee         string `json:"fee"`        // Total for matched amount
	HoldingDays int    `json:"holding_days"`
}

// FifoMoveLotsRequest is the input for fifo_move_lots.
type FifoMoveLotsRequest struct {
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Asset       string `json:"asset"`
	Amount      string `json:"amount"`
}

// FifoMoveLotsResponse is returned by the host from fifo_move_lots.
type FifoMoveLotsResponse struct {
	Result      bool   `json:"result"`
	MovedCount  int    `json:"moved_count"`
	MovedAmount string `json:"moved_amount"` // Actual amount moved (may be less than requested)
	Shortfall   string `json:"shortfall"`    // Requested amount that could not be moved (zero if fully satisfied)
	Error       string `json:"error,omitempty"`
}

// FifoGetBalanceRequest is the input for fifo_get_balance.
type FifoGetBalanceRequest struct {
	Account string `json:"account"`
	Asset   string `json:"asset"`
}

// FifoGetBalanceResponse is returned by the host from fifo_get_balance.
type FifoGetBalanceResponse struct {
	Result  bool   `json:"result"`
	Balance string `json:"balance"`
	Error   string `json:"error,omitempty"`
}

// FifoStateRequest is the input for fifo_save_state and fifo_restore_state.
type FifoStateRequest struct {
	Key string `json:"key"`
}

// FifoStateResponse is returned by fifo_restore_state.
type FifoStateResponse struct {
	Result bool   `json:"result"`
	Found  bool   `json:"found"`
	Error  string `json:"error,omitempty"`
}
