package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"html/template"
	"log"
	"monte-carlo-assignment/calculations"
	"monte-carlo-assignment/handlers"
	"monte-carlo-assignment/ingestion"
	"monte-carlo-assignment/market_data"
	"monte-carlo-assignment/storage"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Setting up logger
	atom := zap.NewAtomicLevel()
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	atom.SetLevel(zap.DebugLevel)

	defer logger.Sync()

	// Setting timeout 30 seconds for all APIs
	ctx := context.Background()
	ctx, cancelCtx := context.WithTimeout(ctx, 30*time.Second)
	defer cancelCtx()

	// Setting up DB
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&storage.QuotePrice{}, &storage.QuoteRank{})

	// Setting up ingestion
	ingestionClient := market_data.NewCryptowatClient(logger)
	quotePriceRepo := storage.NewQuotePriceRepository(logger, db)
	quoteRankRepo := storage.NewQuoteRankRepository(logger, db)
	quotePriceFetcher := ingestion.NewQuotePriceFetcher(logger, ingestionClient, quotePriceRepo)
	rankCalculator := calculations.NewRankCalculator(logger, quotePriceRepo, quoteRankRepo)

	// Setting cron job
	logger.Info("Create new cron")
	c := cron.New()
	c.AddFunc("*/1 * * * *", quotePriceFetcher.FetchQuotePrice)
	c.AddFunc("*/1 * * * *", rankCalculator.CalculateRanks)

	// Start cron with one scheduled job
	logger.Info("Start cron")
	c.Start()

	// Setting up API handlers
	http.HandleFunc("/observability", renderPage)
	quotePriceHandler := handlers.NewQuotePriceHandler(quotePriceRepo, logger)
	http.HandleFunc(handlers.QuotePricePath, quotePriceHandler.ServeRequest)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/index.html")
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
