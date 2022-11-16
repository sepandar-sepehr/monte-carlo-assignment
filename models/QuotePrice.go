package models

import (
	"gorm.io/gorm"
	"time"
)

type QuotePrice struct {
	gorm.Model
	Exchange  string
	Pair      string
	Price     int
	FetchedAt time.Time
}
