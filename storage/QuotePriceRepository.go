package storage

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"monte-carlo-assignment/storage/models"
)

type QuotePriceRepository struct {
	log *zap.Logger
	db  *gorm.DB
}

func NewQuotePriceRepository(log *zap.Logger, db *gorm.DB) QuotePriceRepository {
	return QuotePriceRepository{
		log,
		db,
	}
}

func (r *QuotePriceRepository) StorePrice(quotePrice models.QuotePrice) {
	r.db.Create(&quotePrice)
}

func (r *QuotePriceRepository) Get24hPrice(quotePrice models.QuotePrice) ([]models.QuotePrice, error) {
	minus24hTime := quotePrice.FetchedAt.AddDate(0, 0, -1)
	var prices []models.QuotePrice
	result := r.db.Where("fetched_at >= ? AND fetched_at < ?AND from_symbol = ? AND to_symbol = ? AND exchange = ?",
		minus24hTime, quotePrice.FetchedAt, quotePrice.FromSymbol, quotePrice.ToSymbol, quotePrice.Exchange).Find(&prices)
	if result.Error != nil {
		return nil, errors.New("failed to fetch data")
	}

	return prices, nil
}
