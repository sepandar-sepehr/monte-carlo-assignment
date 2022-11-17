package market_data

type Client interface {
	GetQuotePrice(input GetQuoteInput) (*GetQuoteOutput, error)
}
