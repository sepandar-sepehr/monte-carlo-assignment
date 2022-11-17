package storage

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"monte-carlo-assignment/storage/models"
	"time"
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

func (r *QuotePriceRepository) StorePrice(exchange string, fromSymbol string, toSymbol string, price float32, fetchTime time.Time) {
	quotePrice := models.QuotePrice{
		Exchange:   exchange,
		FromSymbol: fromSymbol,
		ToSymbol:   toSymbol,
		Price:      price,
		FetchedAt:  fetchTime,
	}
	r.db.Create(&quotePrice)
}
