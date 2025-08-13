package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ExchangeResponse struct {
	Date         string              `json:"Date"`
	PreviousDate string              `json:"PreviousDate"`
	PreviousURL  string              `json:"PreviousURL"`
	Timestamp    string              `json:"Timestamp"`
	Valute       map[string]Currency `json:"Valute"`
}

type Currency struct {
	ID       string  `json:"ID"`
	NumCode  string  `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Nominal  int     `json:"Nominal"`
	Name     string  `json:"Name"`
	Value    float64 `json:"Value"`
	Previous float64 `json:"Previous"`
}

func GetExchangeRate(currencyCode string) (string, error) {
	currencyCode = strings.ToUpper(strings.TrimSpace(currencyCode))

	if currencyCode == "" {
		return "", fmt.Errorf("укажите код валюты")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		return "", fmt.Errorf("ошибка соединения с сервисом курсов валют")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ошибка сервиса курсов валют (код %d)", resp.StatusCode)
	}

	var data ExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("ошибка обработки данных курсов валют")
	}

	currency, exists := data.Valute[currencyCode]
	if !exists {
		return getAvailableCurrencies(data.Valute), nil
	}

	return formatExchangeRate(currency), nil
}

func formatExchangeRate(currency Currency) string {
	change := currency.Value - currency.Previous
	changeText := "без изменений"

	if change > 0 {
		changeText = fmt.Sprintf("рост на %.4f ₽", change)
	} else if change < 0 {
		changeText = fmt.Sprintf("падение на %.4f ₽", -change)
	}

	nominalText := ""
	if currency.Nominal > 1 {
		nominalText = fmt.Sprintf(" (за %d %s)", currency.Nominal, currency.CharCode)
	}

	return fmt.Sprintf(`<b>Курс валюты %s - %s</b>

<b>Текущий курс:</b> %.4f ₽%s
<b>Предыдущий курс:</b> %.4f ₽
<b>Изменение:</b> %s

<i>Данные Центрального банка РФ</i>`,
		currency.CharCode,
		currency.Name,
		currency.Value,
		nominalText,
		currency.Previous,
		changeText)
}

func getAvailableCurrencies(valute map[string]Currency) string {
	result := fmt.Sprintf("<b>Валюта не найдена</b>\n\n<b>Доступные валюты:</b>\n")

	// Показываем популярные валюты
	popular := []string{"USD", "EUR", "CNY", "GBP", "JPY", "CHF", "TRY", "KZT", "BYN"}

	for _, code := range popular {
		if currency, exists := valute[code]; exists {
			result += fmt.Sprintf("• %s - %s\n", code, currency.Name)
		}
	}

	result += "\n<i>Пример: /exchange USD</i>"
	return result
}
