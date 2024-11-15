package models

import "errors"

// Валюта с её кодом и названием
type Currency struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Запрос на конвертацию валют
type ConversionRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// Validate проверяет корректность запроса на конвертацию
func (r *ConversionRequest) Validate() error {
	if r.From == "" || r.To == "" {
		return errors.New("from and to currencies must be provided")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

// Ответ на запрос конвертации
type ConversionResponse struct {
	Result              float64            `json:"result"`
	Rate                float64            `json:"rate"`
	AvailableCurrencies map[string]float64 `json:"available_currencies"`
}
