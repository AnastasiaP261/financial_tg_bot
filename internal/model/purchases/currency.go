package purchases

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
	"go.uber.org/zap"
)

func (m *Model) getTodayRates(ctx context.Context, year, month, day int) (currency.RateToRUB, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get today rates")
	defer span.Finish()

	var ok bool
	var rates currency.RateToRUB

	ok, rates, err := m.Repo.GetRate(ctx, year, month, day)
	if err != nil {
		return currency.RateToRUB{}, errors.Wrap(err, "Repo.GetRate")
	}
	if !ok {
		rates, err = m.ExchangeRatesModel.GetExchangeRateToRUBFromDate(ctx, year, month, day)
		if err != nil {
			return currency.RateToRUB{}, errors.Wrap(err, "ExchangeRatesModel.GetExchangeRateToRUBFromDate")
		}

		go func() {
			err := m.Repo.AddRate(ctx, year, month, day, rates)
			if err != nil {
				logs.Error(
					"rate has not been added to the database",
					zap.String("date", fmt.Sprintf("%d.%02d.%02d", year, month, day)),
					zap.Any("rates", rates),
				)
			}
		}()
	}

	return rates, nil
}
