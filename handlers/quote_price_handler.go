package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"monte-carlo-assignment/api_types"
	"monte-carlo-assignment/storage"
	"monte-carlo-assignment/storage/models"
	"net/http"
	"time"
)

type QuotePriceHandler struct {
	log                  *zap.Logger
	quotePriceRepository storage.QuotePriceRepository
}

func NewQuotePriceHandler(quotePriceRepository storage.QuotePriceRepository, log *zap.Logger) *QuotePriceHandler {
	return &QuotePriceHandler{
		log,
		quotePriceRepository,
	}
}

func (h *QuotePriceHandler) ServeRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()
	get24hInput := models.QuotePrice{
		Exchange:   queryParams.Get("exchange"),
		FromSymbol: queryParams.Get("from"),
		ToSymbol:   queryParams.Get("to"),
		FetchedAt:  time.Now(),
	}
	repoResponse, err := h.quotePriceRepository.Get24hPrice(get24hInput)
	if err != nil {
		h.log.Error("failed to get data from repo", zap.Error(err))
		http.Error(w, "Could not fetch prices", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(convertResponse(repoResponse))
}

func convertResponse(repoResponse []models.QuotePrice) []api_types.QuotePrice {
	response := make([]api_types.QuotePrice, len(repoResponse))
	for i, repoPrice := range repoResponse {
		response[i] = api_types.QuotePrice{
			Exchange:   repoPrice.Exchange,
			FromSymbol: repoPrice.FromSymbol,
			ToSymbol:   repoPrice.ToSymbol,
			Price:      repoPrice.Price,
			Time:       repoPrice.FetchedAt,
		}
	}
	return response
}