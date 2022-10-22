package purchases

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/normalize"
)

// AddPurchaseReq тело запроса в Repo для добавления траты
type AddPurchaseReq struct {
	UserID     int64
	Sum        float64
	CategoryID uint64
	Date       time.Time

	// коэффициенты валют на момент совершения траты
	USDRatio float64
	CNYRatio float64
	EURRatio float64
}

// CategoryRow тело запроса в Repo для проверки существования категории у пользователя
type CategoryRow struct {
	UserID   int64
	Category string
}

// AddPurchase добавляет трату.
// Если category пустой, трата будет добавлена без категории.
// Если rawDate пустой, для траты будет выставлена текущая дата.
func (m *Model) AddPurchase(ctx context.Context, userID int64, rawSum, category, rawDate string) error {
	fmt.Println("### AddPurchase")

	var (
		sumCurrency float64
		date        time.Time
		err         error
	)

	sumCurrency, err = strconv.ParseFloat(rawSum, 64)
	if err != nil {
		return ErrSummaParsing
	}

	var categoryID uint64
	if category != "" {
		category = strings.ToLower(category)
		categoryID, err = m.Repo.GetCategoryID(ctx, CategoryRow{
			UserID:   userID,
			Category: normalize.Category(category),
		})

		if err != nil {
			return errors.Wrap(err, "Repo.GetCategoryID")
		}
		if categoryID == 0 {
			return ErrCategoryNotExist
		}
	}

	if rawDate != "" {
		date, err = time.Parse("02.01.2006", rawDate)
		if err != nil {
			return ErrDateParsing
		}
	} else {
		date = time.Now()
	}

	rates := m.ExchangeRatesModel.GetExchangeRateToRUB()

	info, err := m.Repo.GetUserInfo(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "Repo.GetUserInfo")
	}

	sum, err := m.toRUB(info.Currency, sumCurrency, rates)
	if err != nil {
		return errors.Wrap(err, "toRUB")
	}

	if err = m.Repo.AddPurchase(ctx, AddPurchaseReq{
		UserID:     userID,
		Sum:        sum,
		CategoryID: categoryID,
		Date:       date,
		CNYRatio:   rates.CNY,
		EURRatio:   rates.EUR,
		USDRatio:   rates.USD,
	}); err != nil {
		return errors.Wrap(err, "Repo.AddPurchase")
	}

	return nil
}
