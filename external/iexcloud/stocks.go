package iexcloud

import (
	"context"
	"encoding/json"
	"strconv"
)

// Docs: https://iexcloud.io/docs/api/#stocks

// KeyStats provides Key Stats for a Stock
// copied from https://github.com/goinvest/iexcloud
type KeyStats struct {
	Name                string  `json:"companyName"`
	MarketCap           int     `json:"marketCap"`
	Week52High          float64 `json:"week52High"`
	Week52Low           float64 `json:"week52Low"`
	Week52Change        float64 `json:"week52Change"`
	SharesOutstanding   int     `json:"sharesOutstanding"`
	Avg30Volume         float64 `json:"avg30Volume"`
	Avg10Volume         float64 `json:"avg10Volume"`
	Float               int     `json:"float"`
	Symbol              string  `json:"symbol"`
	Employees           int     `json:"employees"`
	TTMEPS              float64 `json:"ttmEPS"`
	TTMDividendRate     float64 `json:"ttmDividendRate"`
	DividendYield       float64 `json:"dividendYield"`
	NextDividendDate    Date    `json:"nextDividendDate"`
	ExDividendDate      Date    `json:"exDividendDate"`
	NextEarningsDate    Date    `json:"nextEarningsDate"`
	PERatio             float64 `json:"peRatio"`
	Day200MovingAvg     float64 `json:"day200MovingAvg"`
	Day50MovingAvg      float64 `json:"day50MovingAvg"`
	MaxChangePercent    float64 `json:"maxChangePercent"`
	Year5ChangePercent  float64 `json:"year5ChangePercent"`
	Year2ChangePercent  float64 `json:"year2ChangePercent"`
	Year1ChangePercent  float64 `json:"year1ChangePercent"`
	YTDChangePercent    float64 `json:"ytdChangePercent"`
	Month6ChangePercent float64 `json:"month6ChangePercent"`
	Month3ChangePercent float64 `json:"month3ChangePercent"`
	Month1ChangePercent float64 `json:"month1ChangePercent"`
	Day30ChangePercent  float64 `json:"day30ChangePercent"`
	Day5ChangePercent   float64 `json:"day5ChangePercent"`
}

func (iex *IEX) StocksKeyStats(ctx context.Context, symbol string) (*KeyStats, error) {
	raw, err := iex.get(ctx, endpointStocksKeyStats(symbol, ""))
	if err != nil {
		return nil, err
	}

	var stats KeyStats
	err = json.Unmarshal(raw, &stats)

	return &stats, err
}

func (iex *IEX) StocksPrice(ctx context.Context, symbol string) (float64, error) {
	raw, err := iex.get(ctx, endpointStocksPrice(symbol))
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(string(raw), 10)
}

// OHLC models the open, high, low, close for a stock.
// copied from https://github.com/goinvest/iexcloud
type OHLC struct {
	Open  OpenClose `json:"open"`
	Close OpenClose `json:"close"`
	High  float64   `json:"high"`
	Low   float64   `json:"low"`
}

// OpenClose provides the price and time for either the open or close price of a stock.
// copied from https://github.com/goinvest/iexcloud
type OpenClose struct {
	Price float64 `json:"price"`
	Time  int     `json:"Time"`
}

func (iex *IEX) StocksOHLC(ctx context.Context, symbol string) (*OHLC, error) {
	raw, err := iex.get(ctx, endpointStocksOHLC(symbol))
	if err != nil {
		return nil, err
	}

	var stats OHLC
	err = json.Unmarshal(raw, &stats)

	return &stats, err
}

// Quote models the data returned from the IEX Cloud /quote endpoint.
// copied from https://github.com/goinvest/iexcloud
type Quote struct {
	Symbol                string    `json:"symbol"`
	CompanyName           string    `json:"companyName"`
	CalculationPrice      string    `json:"calculationPrice"`
	Open                  float64   `json:"open"`
	OpenTime              EpochTime `json:"openTime"`
	Close                 float64   `json:"close"`
	CloseTime             EpochTime `json:"closeTime"`
	High                  float64   `json:"high"`
	Low                   float64   `json:"low"`
	LatestPrice           float64   `json:"latestPrice"`
	LatestSource          string    `json:"latestSource"`
	LatestTime            string    `json:"latestTime"`
	LatestUpdate          EpochTime `json:"latestUpdate"`
	LatestVolume          int       `json:"latestVolume"`
	IEXRealtimePrice      float64   `json:"iexRealtimePrice"`
	IEXRealtimeSize       int       `json:"iexRealtimeSize"`
	IEXLastUpdated        EpochTime `json:"iexLastUpdated"`
	DelayedPrice          float64   `json:"delayedPrice"`
	DelayedPriceTime      EpochTime `json:"delayedPriceTime"`
	ExtendedPrice         float64   `json:"extendedPrice"`
	ExtendedChange        float64   `json:"extendedChange"`
	ExtendedChangePercent float64   `json:"extendedChangePercent"`
	ExtendedPriceTime     EpochTime `json:"extendedPriceTime"`
	PreviousClose         float64   `json:"previousClose"`
	Change                float64   `json:"change"`
	ChangePercent         float64   `json:"changePercent"`
	IEXMarketPercent      float64   `json:"iexMarketPercent"`
	IEXVolume             int       `json:"iexVolume"`
	AvgTotalVolume        int       `json:"avgTotalVolume"`
	IEXBidPrice           float64   `json:"iexBidPrice"`
	IEXBidSize            int       `json:"iexBidSize"`
	IEXAskPrice           float64   `json:"iexAskPrice"`
	IEXAskSize            int       `json:"iexAskSize"`
	MarketCap             int       `json:"marketCap"`
	Week52High            float64   `json:"week52High"`
	Week52Low             float64   `json:"week52Low"`
	YTDChange             float64   `json:"ytdChange"`
	PERatio               float64   `json:"peRatio"`
}

func (iex *IEX) StocksQuote(ctx context.Context, symbol string) (*Quote, error) {
	raw, err := iex.get(ctx, endpointStocksQuote(symbol, ""))
	if err != nil {
		return nil, err
	}

	var stats Quote
	err = json.Unmarshal(raw, &stats)

	return &stats, err
}

type Logo struct {
	URL string `json:"url"`
}

func (iex *IEX) StocksLogo(ctx context.Context, symbol string) (*Logo, error) {
	raw, err := iex.get(ctx, endpointStocksLogo(symbol))
	if err != nil {
		return nil, err
	}

	var stats Logo
	err = json.Unmarshal(raw, &stats)

	return &stats, err
}

type HistoricalPrice struct {
	Date           Date    `json:"date"`
	Open           float64 `json:"open"`
	Close          float64 `json:"close"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Volume         int     `json:"volume"`
	UOpen          float64 `json:"uOpen"`
	UClose         float64 `json:"uClose"`
	UHigh          float64 `json:"uHigh"`
	ULow           float64 `json:"uLow"`
	UVolume        int     `json:"uVolume"`
	Change         float64 `json:"change"`
	ChangePercent  float64 `json:"changePercent"`
	Label          string  `json:"label"`
	ChangeOverTime float64 `json:"changeOverTime"`
}

func (iex *IEX) StocksHistoricalPrices(
	ctx context.Context,
	symbol string,
	rangeData string,
	date string,
) ([]*HistoricalPrice, error) {
	raw, err := iex.get(ctx, endpointStocksHistoricalPrices(symbol, rangeData, date))
	if err != nil {
		return nil, err
	}

	var stats []*HistoricalPrice
	err = json.Unmarshal(raw, &stats)

	return stats, err
}
