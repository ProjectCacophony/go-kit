package iexcloud

import (
	"context"
	"encoding/json"
	"fmt"
)

type Symbol struct {
	Symbol    string `json:"symbol"`
	Exchange  string `json:"exchange"`
	Name      string `json:"name"`
	Date      Date   `json:"date"`
	Type      string `json:"type"`
	IEXID     string `json:"iexId"`
	Region    string `json:"region"`
	Currency  string `json:"currency"`
	IsEnabled bool   `json:"isEnabled"`
}

func (iex *IEX) RefDataSymbols(ctx context.Context) ([]*Symbol, error) {
	raw, err := iex.get(ctx, endpointRefDataSymbols)
	if err != nil {
		return nil, err
	}

	var symbols []*Symbol
	err = json.Unmarshal(raw, &symbols)

	return symbols, err
}

func (iex *IEX) RefDataSymbolsInternational(ctx context.Context, region string) ([]*Symbol, error) {
	raw, err := iex.get(ctx, endpointRefDataSymbolsInternational(region))
	if err != nil {
		return nil, err
	}

	var symbols []*Symbol
	err = json.Unmarshal(raw, &symbols)

	return symbols, err
}

func (s *Symbol) FormatCurrency(myValue float64) string {
	currencyFormat := map[string]string{
		"USD": "$ %.2f",
		"EUR": "%.2f €",
	}

	return fmt.Sprintf(currencyFormat[s.Currency], myValue)
}
