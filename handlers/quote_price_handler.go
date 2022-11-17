package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"monte-carlo-assignment/api_models"
	"monte-carlo-assignment/storage"
	"net/http"
	"time"
)

type QuotePriceHandler struct {
	log                  *zap.Logger
	quotePriceRepository storage.QuotePriceRepository
	quoteRankRepository  storage.QuoteRankRepository
}

func NewQuotePriceHandler(log *zap.Logger, quotePriceRepository storage.QuotePriceRepository,
	quoteRankRepository storage.QuoteRankRepository) *QuotePriceHandler {
	return &QuotePriceHandler{
		log,
		quotePriceRepository,
		quoteRankRepository,
	}
}

func (h *QuotePriceHandler) ServeRequest(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Serving quote price request")

	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()
	storageInput := storage.QuotePrice{
		Exchange:   queryParams.Get("exchange"),
		FromSymbol: queryParams.Get("from"),
		ToSymbol:   queryParams.Get("to"),
		FetchedAt:  time.Now(),
	}
	last24hPrices, err := h.quotePriceRepository.Get24hPrice(storageInput)
	if err != nil {
		h.log.Error("failed to get data from repo", zap.Error(err))
		http.Error(w, "Could not fetch prices", http.StatusInternalServerError)
		return
	}

	latestRank, err := h.quoteRankRepository.GetLatestRank(storageInput)
	if err != nil {
		h.log.Error("failed to get rank from repo", zap.Error(err))
		http.Error(w, "Could not fetch rank", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(convertResponse(last24hPrices, latestRank))
}

func convertResponse(repoResponse []storage.QuotePrice, rank *storage.QuoteRank) *api_models.QuotePriceResponse {
	convertedPrices := make([]api_models.QuotePrice, len(repoResponse))
	for i, repoPrice := range repoResponse {
		convertedPrices[i] = api_models.QuotePrice{
			Exchange:   repoPrice.Exchange,
			FromSymbol: repoPrice.FromSymbol,
			ToSymbol:   repoPrice.ToSymbol,
			Price:      repoPrice.Price,
			Time:       repoPrice.FetchedAt,
		}
	}
	if rank != nil {
		return &api_models.QuotePriceResponse{
			Quotes:          convertedPrices,
			RankNumerator:   rank.RankNumerator,
			RankDenominator: rank.RankDenominator,
		}
	} else {
		return &api_models.QuotePriceResponse{
			Quotes: convertedPrices,
		}
	}
}
