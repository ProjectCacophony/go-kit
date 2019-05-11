package iexcloud

import (
	"fmt"
)

const (
	apiBase = "https://cloud.iexapis.com/v1"

	endpointRefDataSymbols = "/ref-data/symbols"
)

func endpointStocksKeyStats(symbol, stat string) string {
	return fmt.Sprintf("/stock/%s/stats/%s", symbol, stat)
}

func endpointStocksPrice(symbol string) string {
	return fmt.Sprintf("/stock/%s/price", symbol)
}

func endpointStocksOHLC(symbol string) string {
	return fmt.Sprintf("/stock/%s/ohlc", symbol)
}

func endpointStocksQuote(symbol, field string) string {
	return fmt.Sprintf("/stock/%s/quote/%s", symbol, field)
}

func endpointStocksLogo(symbol string) string {
	return fmt.Sprintf("/stock/%s/logo", symbol)
}

func endpointStocksHistoricalPrices(symbol, dataRange, date string) string {
	return fmt.Sprintf("/stock/%s/chart/%s/%s", symbol, dataRange, date)
}

func endpointRefDataSymbolsInternational(region string) string {
	return fmt.Sprintf("/ref-data/region/%s/symbols", region)
}

func (iex *IEX) fullURL(endpoint string) string {
	return apiBase + endpoint + "?token=" + iex.apiKey
}
