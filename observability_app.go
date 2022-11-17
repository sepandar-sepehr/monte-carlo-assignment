package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"html/template"
	"log"
	"monte-carlo-assignment/ingestion"
	"monte-carlo-assignment/market_data"
	"monte-carlo-assignment/storage"
	"monte-carlo-assignment/storage/models"
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
	db.AutoMigrate(&models.QuotePrice{})

	// Setting up ingestion
	ingestionClient := market_data.NewCryptowatClient(logger)
	quotePriceRepo := storage.NewQuotePriceRepository(logger, db)
	quotePriceFetcher := ingestion.NewQuotePriceFetcher(logger, ingestionClient, quotePriceRepo)

	// Setting cron job
	logger.Info("Create new cron")
	c := cron.New()
	c.AddFunc("*/1 * * * *", quotePriceFetcher.FetchQuotePrice)

	// Start cron with one scheduled job
	logger.Info("Start cron")
	c.Start()

	// Setting up web handler
	http.HandleFunc("/observability", renderPage)
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
