package ingestion

import (
	"go.uber.org/zap"
	"monte-carlo-assignment/market_data"
)

type QuotePriceFetcher struct {
	log          *zap.Logger
	marketClient market_data.Client
}

func NewQuotePriceFetcher(log *zap.Logger, client market_data.Client) *QuotePriceFetcher {
	return &QuotePriceFetcher{
		log,
		client,
	}
}
func (f *QuotePriceFetcher) FetchQuotePrice() {
	exchange := "coinbase-pro"
	fromSymbol := "btc"
	toSymbol := "eur"
	input := market_data.GetQuoteInput{
		Exchange:   exchange,
		SymbolFrom: fromSymbol,
		SymbolTo:   toSymbol,
	}
	getQuoteOutput, err := f.marketClient.GetQuotePrice(input)
	if err != nil {
		f.log.Error("failed to fetch price",
			zap.String("exchange", exchange),
			zap.String("fromSymbol", fromSymbol),
			zap.String("toSymbol", toSymbol))
	}
	f.log.Info("fetched price",
		zap.Float32("price", getQuoteOutput.Price),
		zap.String("exchange", exchange),
		zap.String("fromSymbol", fromSymbol),
		zap.String("toSymbol", toSymbol),
	)
}
