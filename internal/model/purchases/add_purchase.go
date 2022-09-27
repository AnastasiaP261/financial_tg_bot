package purchases

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

// AddPurchaseReq тело запроса в Repo для добавления траты
type AddPurchaseReq struct {
	User     string
	Sum      float64
	Category string
	Date     time.Time
}

// CategoryExistReq тело запроса в Repo для проверки существования категории у пользователя
type CategoryExistReq struct {
	User     string
	Category string
}

// AddPurchase добавляет трату.
// Если category пустой, трата будет добавлена без категории.
// Если rawDate пустой, для траты будет выставлено время добавления.
func (m *Model) AddPurchase(user, rawSum, category, rawDate string) error {
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
		categoryExist = m.Repo.CategoryExist(CategoryExistReq{
			User:     user,
			Category: category,
		})
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
		User:     user,
		Sum:      sum,
		Category: category,
		Date:     date,
	}); err != nil {
		return errors.Wrap(err, "Repo.AddPurchase")
	}

	return nil
}
