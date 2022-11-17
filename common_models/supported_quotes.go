package common_models

type quotePair struct {
	FromSymbol string
	ToSymbol   string
}

const SupportedExchange = "coinbase-pro"

var SupportedQuotes = []quotePair{
	{
		FromSymbol: "btc",
		ToSymbol:   "eur",
	},
	{
		FromSymbol: "bnt",
		ToSymbol:   "eur",
	},
	{
		FromSymbol: "eth",
		ToSymbol:   "eur",
	},
}
