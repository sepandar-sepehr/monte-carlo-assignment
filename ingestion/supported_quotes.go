package ingestion

type quotePair struct {
	FromSymbol string
	ToSymbol   string
}

var supportedQuotes = []quotePair{
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
