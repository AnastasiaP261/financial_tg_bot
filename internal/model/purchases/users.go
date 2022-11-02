package purchases

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/normalize"
)

type User struct {
	UserID     int64
	Currency   Currency // выбранная пользователем валюта
	Categories []int64
	Limit      float64
}

func (m *Model) ChangeUserCurrency(ctx context.Context, userID int64, currency Currency) error {
	err := m.Repo.ChangeCurrency(ctx, userID, currency)
	if err != nil {
		return errors.Wrap(err, "Repo.ChangeCurrency")
	}
	return nil
}

func (m *Model) ChangeUserLimit(ctx context.Context, userID int64, rawLimit string) error {
	limitCurrency, err := strconv.ParseFloat(rawLimit, 64)
	if err != nil {
		return ErrLimitParsing
	}

	info, err := m.Repo.GetUserInfo(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "Repo.GetUserInfo")
	}

	y, month, d := time.Now().Date()
	rates, err := m.getTodayRates(ctx, y, int(month), d)
	if err != nil {
		return errors.Wrap(err, "getTodayRates")
	}

	limit, err := m.toRUB(info.Currency, limitCurrency, rates)
	if err != nil {
		return errors.Wrap(err, "limit to rub")
	}

	err = m.Repo.ChangeUserLimit(ctx, userID, limit)
	if err != nil {
		return errors.Wrap(err, "Repo.ChangeCurrency")
	}
	return nil
}

// AddCategoryToUser добавить новую категорию пользователю
func (m *Model) AddCategoryToUser(ctx context.Context, userID int64, category string) error {
	if err := m.Repo.AddCategoryToUser(ctx, userID, normalize.Category(category)); err != nil {
		return errors.Wrap(err, "Repo.AddCategoryToUser")
	}
	return nil
}

func (m *Model) GetUserCategories(ctx context.Context, userID int64) ([]string, error) {
	res, err := m.Repo.GetUserCategories(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "Repo.AddCategoryToUser")
	}
	return res, nil
}
