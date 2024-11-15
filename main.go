package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Jorik887/serviceConvertCoin/handlers"
	"github.com/Jorik887/serviceConvertCoin/services"
)

func main() {
	// Запуск горутины для обновления курсов валют
	go func() {
		for {
			// Получение курсов валют из внешнего API
			_, err := services.GetCurrencyRates()
			if err != nil {
				log.Printf("Ошибка при получении курсов валют: %v", err)
				time.Sleep(1 * time.Minute)
				continue
			}

			log.Println("Курсы валют обновлены")
			time.Sleep(10 * time.Minute) // Обновление каждые 10 минут
		}
	}()

	// Настройка маршрутов
	http.HandleFunc("/convert", handlers.ConvertCurrency)

	// Запуск сервера
	port := ":8080"
	log.Printf("Сервер запущен на порту %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
