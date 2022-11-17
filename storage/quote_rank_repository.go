package storage

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type QuoteRankRepository struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewQuoteRankRepository(log *zap.Logger, db *gorm.DB) QuoteRankRepository {
	return QuoteRankRepository{
		log,
		db,
	}
}

func (r *QuoteRankRepository) StoreRank(quoteRank QuoteRank) {
	r.db.Create(&quoteRank)
}

// GetLatestRank gives 5min grace period in case rank calculation is delayed for some reason, so we can use the latest
func (r *QuoteRankRepository) GetLatestRank(quotePrice QuotePrice) (*QuoteRank, error) {
	minus5mTime := quotePrice.FetchedAt.Add(-time.Minute * 5)
	var quoteRank QuoteRank
	result := r.db.Where("calculated_at >= ? AND calculated_at < ? AND from_symbol = ? AND to_symbol = ? AND exchange = ?",
		minus5mTime, quotePrice.FetchedAt, quotePrice.FromSymbol, quotePrice.ToSymbol, quotePrice.Exchange).
		Order("calculated_at desc").
		First(&quoteRank)
	if result.Error != nil {
		r.log.Error(result.Error.Error())
		return nil, errors.New("failed to fetch data")
	}

	return &quoteRank, nil
}
