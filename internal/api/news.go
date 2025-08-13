package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NewsResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Source struct {
		ID   *string `json:"id"`
		Name string  `json:"name"`
	} `json:"source"`
	Author      *string `json:"author"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	URL         string  `json:"url"`
	URLToImage  *string `json:"urlToImage"`
	PublishedAt string  `json:"publishedAt"`
	Content     *string `json:"content"`
}

func GetNews(apiKey string) (string, error) {
	if apiKey == "" {
		return getNewsStub(), nil
	}

	// Используем top-headlines для главных новостей России
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=ru&pageSize=5&apiKey=%s", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка соединения с сервисом новостей")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return "", fmt.Errorf("неверный API ключ NewsAPI")
	}

	if resp.StatusCode == 429 {
		return "", fmt.Errorf("превышен лимит запросов к API новостей")
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ошибка сервиса новостей (код %d)", resp.StatusCode)
	}

	var news NewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return "", fmt.Errorf("ошибка обработки данных новостей")
	}

	if news.Status != "ok" {
		return "", fmt.Errorf("ошибка получения новостей")
	}

	if len(news.Articles) == 0 {
		return "На данный момент новостей нет", nil
	}

	return formatNews(news.Articles), nil
}

func formatNews(articles []Article) string {
	result := "<b>Главные новости дня</b>\n\n"

	for i, article := range articles {
		if i >= 5 { // максимум 5 новостей
			break
		}

		title := article.Title
		if len(title) > 100 {
			title = title[:97] + "..."
		}

		source := "Неизвестный источник"
		if article.Source.Name != "" {
			source = article.Source.Name
		}

		result += fmt.Sprintf("<b>%d. %s</b>\n", i+1, title)

		// Добавляем описание если есть
		if article.Description != nil && *article.Description != "" {
			description := *article.Description
			if len(description) > 150 {
				description = description[:147] + "..."
			}
			result += fmt.Sprintf("%s\n", description)
		}

		result += fmt.Sprintf("<i>Источник: %s</i>\n\n", source)
	}

	result += "<i>Данные предоставлены NewsAPI</i>"
	return result
}

func getNewsStub() string {
	return `<b>Главные новости дня (демо-режим)</b>

<b>1. Российские IT-специалисты показывают рост зарплат</b>
Средняя зарплата разработчиков выросла на 15% за последний год согласно исследованию рекрутингового агентства.
<i>Источник: РБК</i>

<b>2. Удаленная работа становится стандартом для IT-сферы</b>
85% российских IT-компаний готовы предоставить сотрудникам возможность полностью удаленной работы.
<i>Источник: Ведомости</i>

<b>3. Искусственный интеллект меняет рынок труда</b>
Появляются новые профессии связанные с разработкой и внедрением ИИ-решений в российских компаниях.
<i>Источник: Коммерсант</i>

<b>4. Рост спроса на Go-разработчиков</b>
Язык программирования Go показывает увеличение вакансий на 40% по сравнению с прошлым годом.
<i>Источник: HeadHunter</i>

<b>5. Новые меры поддержки IT-отрасли</b>
Правительство анонсировало дополнительные льготы для IT-компаний и специалистов.
<i>Источник: ТАСС</i>

<i>Для получения актуальных новостей настройте NEWS_API_KEY</i>`
}
