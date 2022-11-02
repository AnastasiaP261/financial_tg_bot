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
	ID       uint64
	Category string
}

type ExpensesAndLimit struct {
	Limit         float64  // установленный пользователем лимит (в выбранной валюте)
	Expenses      float64  // сколько он уже потратил за месяц (если лимит установлен) (в выбранной валюте)
	Currency      Currency // выбранная валюта
	LimitExceeded bool     // превышен ли лимит
}

// AddPurchase добавляет трату.
// Если category пустой, трата будет добавлена без категории.
// Если rawDate пустой, для траты будет выставлена текущая дата.
func (m *Model) AddPurchase(ctx context.Context, userID int64, rawSum, category, rawDate string) (ExpensesAndLimit, error) {
	var (
		sumCurrency float64
		date        time.Time
		err         error
	)

	sumCurrency, err = strconv.ParseFloat(rawSum, 64)
	if err != nil {
		return ExpensesAndLimit{}, ErrSummaParsing
	}

	// получаем id категории которую выбрал пользователь и проверяем что такая категория существует
	var categoryID uint64
	if category != "" {
		category = strings.ToLower(category)
		categoryID, err = m.Repo.GetCategoryID(ctx, normalize.Category(category))
		if err != nil {
			return ExpensesAndLimit{}, errors.Wrap(err, "Repo.GetCategoryID")
		}
		if categoryID == 0 {
			return ExpensesAndLimit{}, ErrCategoryNotExist
		}
	} else {
		categoryID = 1
	}

	fmt.Println("### categoryID", categoryID)
	// проверяем, создана ли такая категория у юзера
	has, err := m.Repo.UserHasCategory(ctx, userID, categoryID)
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "Repo.UserHasCategory")
	}
	if !has {
		return ExpensesAndLimit{}, ErrUserHasntCategory
	}

	// переводим сумму траты которую он ввел в рубли
	rates := RateToRUB{}
	if rawDate != "" {
		date, err = time.Parse("02.01.2006", rawDate)
		if err != nil {
			return ExpensesAndLimit{}, errors.Wrap(ErrInvalidDate, "parsing err")
		}

		day, month, year, err := RawDateToYMD(rawDate)
		if err != nil {
			return ExpensesAndLimit{}, errors.Wrap(err, "RawDateToYMD")
		}

		rates, err = m.getTodayRates(ctx, year, month, day)
		if err != nil {
			return ExpensesAndLimit{}, errors.Wrap(err, "getTodayRates")
		}
	} else {
		date = time.Now()
		rates = m.ExchangeRatesModel.GetExchangeRateToRUB()
	}

	info, err := m.Repo.GetUserInfo(ctx, userID)
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "Repo.GetUserInfo")
	}

	sumRUB, err := m.toRUB(info.Currency, sumCurrency, rates)
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "toRUB")
	}

	fmt.Println("### userCur", info.Currency)
	// определяем превышен ли лимит и сколько потрачено за этот календарный месяц
	expAndLim, err := m.getExpensesAndLimit(ctx, userID, info.Currency, info.Limit, sumCurrency, rates)
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "getExpensesAndLimit")
	}

	if err = m.Repo.AddPurchase(ctx, AddPurchaseReq{
		UserID:     userID,
		Sum:        sumRUB,
		CategoryID: categoryID,
		Date:       date,
		CNYRatio:   rates.CNY,
		EURRatio:   rates.EUR,
		USDRatio:   rates.USD,
	}); err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "Repo.AddPurchase")
	}

	return expAndLim, nil
}

func RawDateToYMD(rawDate string) (year, month, day int, err error) {
	t := strings.Split(rawDate, ".")
	if len(t) != 3 {
		return 0, 0, 0, ErrInvalidDate
	}
	y, m, d := t[0], t[1], t[2]

	y1, err := strconv.ParseInt(y, 10, 64)
	if err != nil {
		return 0, 0, 0, ErrInvalidDate
	}
	m1, err := strconv.ParseInt(m, 10, 64)
	if err != nil {
		return 0, 0, 0, ErrInvalidDate
	}
	d1, err := strconv.ParseInt(d, 10, 64)
	if err != nil {
		return 0, 0, 0, ErrInvalidDate
	}

	return int(y1), int(m1), int(d1), nil
}

func (m *Model) getExpensesAndLimit(ctx context.Context, userID int64, userCurrency Currency, userLimit float64, purchaseSum float64, rates RateToRUB) (ExpensesAndLimit, error) {
	expAndLim := ExpensesAndLimit{Limit: userLimit}
	// получаем траты юзера за текущий месяц в рублях
	expRUB, err := m.Repo.GetUserPurchasesSumFromMonth(ctx, userID, time.Now())
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "Repo.GetUserPurchasesSumFromMonth")
	}

	limit, err := m.rubToCurrentCurrency(userCurrency, userLimit, rates)
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "getting limit from rubToCurrentCurrency")
	}

	expenses, err := m.rubToCurrentCurrency(userCurrency, expRUB, rates)
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "getting expenses from rubToCurrentCurrency")
	}
	expenses += purchaseSum

	var limitExceeded bool
	if expAndLim.Limit != -1 && expenses > limit {
		limitExceeded = true
	}

	return ExpensesAndLimit{
		Limit:         limit,
		Expenses:      expenses,
		Currency:      userCurrency,
		LimitExceeded: limitExceeded,
	}, nil
}
