package purchases

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/normalize"
)

// AddPurchaseReq тело запроса в Repo для добавления траты
type AddPurchaseReq struct {
	UserID   int64
	Sum      float64
	Category string
	Date     time.Time
}

// CategoryRow тело запроса в Repo для проверки существования категории у пользователя
type CategoryRow struct {
	UserID   int64
	Category string
}

// AddPurchase добавляет трату.
// Если category пустой, трата будет добавлена без категории.
// Если rawDate пустой, для траты будет выставлена текущая дата.
func (m *Model) AddPurchase(userID int64, rawSum, category, rawDate string) error {
	var (
		sum           float64
		categoryExist bool
		date          time.Time
		err           error
	)

	sum, err = strconv.ParseFloat(rawSum, 64)
	if err != nil {
		return ErrSummaParsing
	}

	if category != "" {
		category = strings.ToLower(category)
		categoryExist, err = m.Repo.CategoryExist(CategoryRow{
			UserID:   userID,
			Category: normalize.Category(category),
		})
		if err != nil {
			return errors.Wrap(err, "Repo.CategoryExist")
		}
		if !categoryExist {
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

	if err = m.Repo.AddPurchase(AddPurchaseReq{
		UserID:   userID,
		Sum:      sum,
		Category: category,
		Date:     date,
	}); err != nil {
		return errors.Wrap(err, "Repo.AddPurchase")
	}

	return nil
}
