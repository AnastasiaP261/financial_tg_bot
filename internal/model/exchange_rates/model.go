package exchange_rates

import (
	"context"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"

	"github.com/pkg/errors"
)

// Client клиент от которого получаем информацию о курсах валют
type Client interface {
	GetExchangeRateToRUB() map[string]float64
	GetExchangeRateToRUBFromDate(ctx context.Context, y, m, d int) (map[string]float64, error)
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

// GetExchangeRateToRUB получить последние валютные курсы для USD, EUR и CNY к RUB
func (c *Model) GetExchangeRateToRUB() currency.RateToRUB {
	rates := c.client.GetExchangeRateToRUB()

	var resp currency.RateToRUB
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

// GetExchangeRateToRUBFromDate получить валютные курсы для USD, EUR и CNY к RUB на определенную дату
func (c *Model) GetExchangeRateToRUBFromDate(ctx context.Context, y, m, d int) (currency.RateToRUB, error) {
	rates, err := c.client.GetExchangeRateToRUBFromDate(ctx, y, m, d)
	if err != nil {
		return currency.RateToRUB{}, errors.Wrap(err, "client.GetExchangeRateToRUBFromDate")
	}

	var resp currency.RateToRUB
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

	return resp, nil
}
