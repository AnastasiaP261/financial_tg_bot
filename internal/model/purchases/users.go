package purchases

import (
	"context"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/normalize"
)

type User struct {
	UserID     int64
	Currency   currency.Currency // выбранная пользователем валюта
	Categories []int64
	Limit      float64
}

func (m *Model) ChangeUserCurrency(ctx context.Context, userID int64, currency currency.Currency) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "change user currency")
	defer span.Finish()

	err := m.Repo.ChangeCurrency(ctx, userID, currency)
	if err != nil {
		return errors.Wrap(err, "repo.ChangeCurrency")
	}
	return nil
}

func (m *Model) ChangeUserLimit(ctx context.Context, userID int64, rawLimit string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "change user limit")
	defer span.Finish()

	limitCurrency, err := strconv.ParseFloat(rawLimit, 64)
	if err != nil {
		return ErrLimitParsing
	}

	info, err := m.Repo.GetUserInfo(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "repo.GetUserInfo")
	}

	y, month, d := time.Now().Date()
	rates, err := m.getTodayRates(ctx, y, int(month), d)
	if err != nil {
		return errors.Wrap(err, "getTodayRates")
	}

	limit, err := currency.ToRUB(info.Currency, limitCurrency, rates)
	if err != nil {
		return errors.Wrap(err, "limit to rub")
	}

	err = m.Repo.ChangeUserLimit(ctx, userID, limit)
	if err != nil {
		return errors.Wrap(err, "repo.ChangeCurrency")
	}
	return nil
}

// AddCategoryToUser добавить новую категорию пользователю
func (m *Model) AddCategoryToUser(ctx context.Context, userID int64, category string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "add category to user")
	defer span.Finish()

	if err := m.Repo.AddCategoryToUser(ctx, userID, normalize.Category(category)); err != nil {
		return errors.Wrap(err, "repo.AddCategoryToUser")
	}
	return nil
}

func (m *Model) GetUserCategories(ctx context.Context, userID int64) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get user categories")
	defer span.Finish()

	res, err := m.Repo.GetUserCategories(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "repo.AddCategoryToUser")
	}
	return res, nil
}
