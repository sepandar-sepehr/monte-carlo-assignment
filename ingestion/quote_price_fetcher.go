package ingestion

import (
	"go.uber.org/zap"
)

type QuotePriceFetcher struct {
	log *zap.Logger
}

func NewQuotePriceFetcher(log *zap.Logger) *QuotePriceFetcher {
	return &QuotePriceFetcher{
		log,
	}
}
func (f *QuotePriceFetcher) FetchQuotePrice() {
	f.log.Info("[Job 1]Every minute job")
}
