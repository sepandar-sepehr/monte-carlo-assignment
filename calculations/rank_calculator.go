package calculations

import (
	stats "github.com/montanaflynn/stats"
	"go.uber.org/zap"
	"monte-carlo-assignment/common_models"
	"monte-carlo-assignment/storage"
	"time"
)

type RankCalculator struct {
	log                  *zap.Logger
	quotePriceRepository storage.QuotePriceRepository
	quoteRankRepository  storage.QuoteRankRepository
}

func NewRankCalculator(log *zap.Logger, quotePriceRepository storage.QuotePriceRepository,
	quoteRankRepository storage.QuoteRankRepository) *RankCalculator {
	return &RankCalculator{
		log,
		quotePriceRepository,
		quoteRankRepository,
	}
}

func (c *RankCalculator) CalculateRanks() {
	c.log.Info("calculating ranks")

	quotesCount := len(common_models.SupportedQuotes)

	datasets := make([]stats.Float64Data, quotesCount)
	for i, quotePair := range common_models.SupportedQuotes {
		fromSymbol := quotePair.FromSymbol
		toSymbol := quotePair.ToSymbol
		get24hInput := storage.QuotePrice{
			Exchange:   common_models.SupportedExchange,
			FromSymbol: fromSymbol,
			ToSymbol:   toSymbol,
			FetchedAt:  time.Now(),
		}
		repoResponse, err := c.quotePriceRepository.Get24hPrice(get24hInput)
		if err != nil {
			c.log.Error("failed to get data from repo", zap.Error(err))
			return
		}
		datasets[i] = make([]float64, len(repoResponse))
		quoteDataset := datasets[i]
		for j, repoPrice := range repoResponse {
			quoteDataset[j] = float64(repoPrice.Price)
		}
	}

	standardDevs := make(stats.Float64Data, quotesCount)
	for i, quote := range common_models.SupportedQuotes {
		standardDev, err := stats.StandardDeviation(datasets[i])
		if err != nil {
			c.log.Error("failed to calculate standard deviation", zap.String("quote", quote.FromSymbol+quote.ToSymbol))
			return
		}
		standardDevs[i] = standardDev
	}

	for i, quotePair := range common_models.SupportedQuotes {
		fromSymbol := quotePair.FromSymbol
		toSymbol := quotePair.ToSymbol
		rank := c.calcRank(standardDevs[i], standardDevs)
		quoteRank := storage.QuoteRank{
			Exchange:        common_models.SupportedExchange,
			FromSymbol:      fromSymbol,
			ToSymbol:        toSymbol,
			RankNumerator:   rank,
			RankDenominator: quotesCount,
			CalculatedAt:    time.Now(),
		}
		c.log.Info("storing ranks", zap.String("from", fromSymbol), zap.String("to", toSymbol),
			zap.Int("numerator", rank), zap.Int("denominator", quotesCount))
		c.quoteRankRepository.StoreRank(quoteRank)
	}
}

func (c *RankCalculator) calcRank(selectedStdDev float64, standardDevs []float64) int {
	rank := 0
	for _, comparingStdDev := range standardDevs {
		if comparingStdDev <= selectedStdDev {
			rank += 1
		}
	}
	return rank
}
