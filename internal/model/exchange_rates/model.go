package exchange_rates

import (
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

// Client клиент от которого получаем информацию о курсах валют
type Client interface {
	GetExchangeRateToRUB() map[string]float64
}

// Model модель валютных курсов
type Model struct {
	client Client
}

// New создает инстанс модели
func New(client Client) *Model {
	return &Model{
		client: client,
	}
}

// GetExchangeRateToRUB получить валютные курсы для USD, EUR и CNY к RUB
func (c *Model) GetExchangeRateToRUB() purchases.RateToRUB {
	rates := c.client.GetExchangeRateToRUB()

	var resp purchases.RateToRUB
	for k, v := range rates {
		if k == "EUR" {
			resp.EUR = v
			continue
		}
		if k == "USD" {
			resp.USD = v
			continue
		}
		if k == "CNY" {
			resp.CNY = v
		}
	}

	return resp
}
