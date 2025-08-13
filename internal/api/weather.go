package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int64 `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type ErrorResponse struct {
	Cod     int    `json:"cod"`
	Message string `json:"message"`
}

func GetWeather(city, apiKey string) (string, error) {
	if apiKey == "" {
		return getWeatherStub(city), nil
	}

	baseURL := "https://api.openweathermap.org/data/2.5/weather"
	params := url.Values{}
	params.Add("q", city)
	params.Add("appid", apiKey)
	params.Add("units", "metric")
	params.Add("lang", "ru")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("ошибка соединения с сервисом погоды")
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode == 404 {
		return "", fmt.Errorf("город '%s' не найден", city)
	}

	if resp.StatusCode == 401 {
		return "", fmt.Errorf("неверный API ключ OpenWeather")
	}

	if resp.StatusCode != 200 {
		// Пытаемся получить детальную ошибку
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return "", fmt.Errorf("ошибка API: %s", errorResp.Message)
		}
		return "", fmt.Errorf("ошибка сервиса погоды (код %d)", resp.StatusCode)
	}

	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return "", fmt.Errorf("ошибка обработки данных о погоде")
	}

	// Проверяем код ответа в JSON
	if weather.Cod != 200 {
		return "", fmt.Errorf("ошибка получения данных о погоде")
	}

	return formatWeather(weather), nil
}

func formatWeather(w WeatherResponse) string {
	temp := int(w.Main.Temp)
	feelsLike := int(w.Main.FeelsLike)
	tempMin := int(w.Main.TempMin)
	tempMax := int(w.Main.TempMax)

	description := "ясно"
	if len(w.Weather) > 0 {
		description = w.Weather[0].Description
	}

	windSpeed := int(w.Wind.Speed)
	visibility := w.Visibility / 1000 // конвертируем в км

	// Давление в мм рт.ст.
	pressureMmHg := int(float64(w.Main.Pressure) * 0.75006)

	return fmt.Sprintf(`%s <b>%s, %s</b>

<b>Температура:</b> %d°C (ощущается как %d°C)
<b>Мин/Макс:</b> %d°C / %d°C
<b>Описание:</b> %s
<b>Влажность:</b> %d%%
<b>Ветер:</b> %d м/с
<b>Давление:</b> %d мм рт.ст.
<b>Видимость:</b> %d км
<b>Облачность:</b> %d%%`,
		w.Name, w.Sys.Country,
		temp, feelsLike, tempMin, tempMax,
		description, w.Main.Humidity, windSpeed,
		pressureMmHg, visibility, w.Clouds.All)
}

func getWeatherStub(city string) string {
	return fmt.Sprintf(`<b>%s (демо-режим)</b>

<b>Температура:</b> 22°C (ощущается как 24°C)
<b>Мин/Макс:</b> 19°C / 25°C
<b>Описание:</b> переменная облачность
<b>Влажность:</b> 65%%
<b>Ветер:</b> 3 м/с
<b>Давление:</b> 760 мм рт.ст.
<b>Видимость:</b> 10 км
<b>Облачность:</b> 40%%

<i>Для реальных данных настройте OPENWEATHER_API_KEY</i>`, city)
}
