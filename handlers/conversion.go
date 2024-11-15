package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Jorik887/serviceConvertCoin/cache"
	"github.com/Jorik887/serviceConvertCoin/models"
	"github.com/Jorik887/serviceConvertCoin/services"
)

var currencyCache = cache.NewCache()

func ConvertCurrency(w http.ResponseWriter, r *http.Request) {
	var request models.ConversionRequest

	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Декодирование JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Приведение валют к верхнему регистру
	request.From = strings.ToUpper(request.From)
	request.To = strings.ToUpper(request.To)

	// Валидация запроса
	if err := request.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение курса валют из кеша
	fromRate, exists := currencyCache.Get(request.From)
	if !exists {
		// Если курс не найден в кеше, то получаем его из API
		rates, err := services.GetCurrencyRates()
		if err != nil {
			http.Error(w, "Failed to get currency rates: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Сохранение курса в кеш
		fromRate, fromExists := rates.Rates[request.From]
		if !fromExists {
			http.Error(w, fmt.Sprintf("Rate not found for currency: %s", request.From), http.StatusBadRequest)
			return
		}
		// Установление времени существования кеша
		currencyCache.Set(request.From, fromRate, 10*time.Minute)
	}

	// Получение курса для конвертации
	rates, err := services.GetCurrencyRates() // Получаем курсы валют из API
	if err != nil {
		http.Error(w, "Failed to get currency rates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	toRate, toExists := rates.Rates[request.To]
	if !toExists {
		http.Error(w, fmt.Sprintf("Rate not found for currency: %s", request.To), http.StatusBadRequest)
		return
	}

	// Конвертация суммы
	baseAmount := request.Amount / fromRate
	result := baseAmount * toRate

	response := models.ConversionResponse{
		Result:              result,
		Rate:                toRate,
		AvailableCurrencies: rates.Rates,
	}

	// Установление заголовка и кода ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Устанавливаем код состояния только один раз

	// Кодирование ответа в JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
