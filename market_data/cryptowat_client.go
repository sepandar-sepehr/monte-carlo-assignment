package market_data

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type CryptowatClient struct {
	log *zap.Logger
}

func NewCryptowatClient(log *zap.Logger) Client {
	return &CryptowatClient{
		log,
	}
}

func (c CryptowatClient) GetQuotePrice(input GetQuoteInput) (*GetQuoteOutput, error) {
	url := fmt.Sprintf("https://api.cryptowat.ch/markets/%s/%s%s/price", input.Exchange, input.SymbolFrom, input.SymbolTo)
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.New("failed to make client call")
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("failed to read client response")
	}

	var clientResponse cwGetQuoteResponse
	err = json.Unmarshal(responseData, &clientResponse)
	if err != nil {
		return nil, errors.New("failed to deserialize client response")
	}

	return &GetQuoteOutput{
		Price: clientResponse.Result.Price,
	}, nil
}
