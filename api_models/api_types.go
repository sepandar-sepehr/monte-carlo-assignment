package api_models

import "time"

type QuotePriceResponse struct {
	Quotes          []QuotePrice
	RankNumerator   int
	RankDenominator int
}

type QuotePrice struct {
	Exchange   string
	FromSymbol string
	ToSymbol   string
	Price      float32
	Time       time.Time
}
