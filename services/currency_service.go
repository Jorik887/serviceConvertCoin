package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CurrencyRate struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

type ApiResponse struct {
	Data map[string]float64 `json:"data"`
}

// GetCurrencyRates получает актуальные курсы валют
func GetCurrencyRates() (CurrencyRate, error) {
	resp, err := http.Get("https://api.freecurrencyapi.com/v1/latest?apikey=fca_live_VteUgeAT9c0kM5mU5s67zxo2m8hjH9jVKAIo9ZZx")
	if err != nil {
		return CurrencyRate{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CurrencyRate{}, fmt.Errorf("failed to get rates, status code: %d", resp.StatusCode)
	}

	var apiResponse ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return CurrencyRate{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Преобразование данных в CurrencyRate
	rates := CurrencyRate{
		Base:  "USD",
		Rates: apiResponse.Data,
	}

	log.Printf("Получены курсы валют: %+v", rates.Rates)
	return rates, nil
}
