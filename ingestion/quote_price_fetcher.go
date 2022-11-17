package ingestion

import (
	"go.uber.org/zap"
	"monte-carlo-assignment/market_data"
	"monte-carlo-assignment/storage"
	"monte-carlo-assignment/storage/models"
	"time"
)

type QuotePriceFetcher struct {
	log                  *zap.Logger
	marketClient         market_data.Client
	quotePriceRepository storage.QuotePriceRepository
}

func NewQuotePriceFetcher(
	log *zap.Logger,
	client market_data.Client,
	quotePriceRepository storage.QuotePriceRepository,
) *QuotePriceFetcher {
	return &QuotePriceFetcher{
		log,
		client,
		quotePriceRepository,
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

	quotePrice := models.QuotePrice{
		Exchange:   exchange,
		FromSymbol: fromSymbol,
		ToSymbol:   toSymbol,
		Price:      getQuoteOutput.Price,
		FetchedAt:  time.Now(),
	}
	f.quotePriceRepository.StorePrice(quotePrice)
}
