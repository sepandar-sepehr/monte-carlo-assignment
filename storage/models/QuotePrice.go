package models

import (
	"gorm.io/gorm"
	"time"
)

type QuotePrice struct {
	gorm.Model
	Exchange   string
	FromSymbol string
	ToSymbol   string
	Price      float32
	FetchedAt  time.Time
}
