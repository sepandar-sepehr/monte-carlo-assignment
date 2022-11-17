package market_data

type GetQuoteInput struct {
	Exchange   string
	SymbolFrom string
	SymbolTo   string
}

type GetQuoteOutput struct {
	Price float32
}

type cwGetQuoteResponse struct {
	Result cwGetQuoteResult `json:"result"`
}

type cwGetQuoteResult struct {
	Price float32 `json:"price"`
}
