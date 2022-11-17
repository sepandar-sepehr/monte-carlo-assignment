package api_types

import "time"

type QuotePriceResponse struct {
	Quotes []QuotePrice
}

type QuotePrice struct {
	Exchange   string
	FromSymbol string
	ToSymbol   string
	Price      float32
	Time       time.Time
}
