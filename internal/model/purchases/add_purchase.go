package purchases

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
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
	Limit         float64           // установленный пользователем лимит (в выбранной валюте)
	Expenses      float64           // сколько он уже потратил за месяц (если лимит установлен) (в выбранной валюте)
	Currency      currency.Currency // выбранная валюта
	LimitExceeded bool              // превышен ли лимит
}

// AddPurchase добавляет трату.
// Если category пустой, трата будет добавлена без категории.
// Если rawDate пустой, для траты будет выставлена текущая дата.
func (m *Model) AddPurchase(ctx context.Context, userID int64, rawSum, category, rawDate string) (ExpensesAndLimit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "add purchase")
	defer span.Finish()

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
			return ExpensesAndLimit{}, errors.Wrap(err, "repo.GetCategoryID")
		}
		if categoryID == 0 {
			return ExpensesAndLimit{}, ErrCategoryNotExist
		}

		// проверяем, создана ли такая категория у юзера
		has, err := m.Repo.UserHasCategory(ctx, userID, categoryID)
		if err != nil {
			return ExpensesAndLimit{}, errors.Wrap(err, "repo.UserHasCategory")
		}
		if !has {
			return ExpensesAndLimit{}, ErrUserHasntCategory
		}
	} else {
		categoryID = 1
	}

	// переводим сумму траты которую он ввел в рубли
	rates := currency.RateToRUB{}
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
		return ExpensesAndLimit{}, errors.Wrap(err, "repo.GetUserInfo")
	}

	sumRUB, err := currency.ToRUB(info.Currency, sumCurrency, rates)
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "toRUB")
	}

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
		return ExpensesAndLimit{}, errors.Wrap(err, "repo.AddPurchase")
	}

	// если не удалить отчет уже устаревший отчет, то при создании отчета нужно будет проверять,
	// что дата последней совершенной траты не свежее чем дата создания отчета. В таком случае
	// использование кеша будет абсолютно неэффективным, так как все равно придется сделать запрос в бд
	m.ReportsStore.Delete(ctx, createKeyForReportsStore(userID)) // nolint: errcheck

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

func (m *Model) getExpensesAndLimit(ctx context.Context, userID int64, userCurrency currency.Currency, userLimit float64, purchaseSum float64, rates currency.RateToRUB) (ExpensesAndLimit, error) {
	expAndLim := ExpensesAndLimit{Limit: userLimit}
	// получаем траты юзера за текущий месяц в рублях
	expRUB, err := m.Repo.GetUserPurchasesSumFromMonth(ctx, userID, time.Now())
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "repo.GetUserPurchasesSumFromMonth")
	}

	limit, err := currency.RubToCurrentCurrency(userCurrency, userLimit, rates)
	if err != nil {
		return ExpensesAndLimit{}, errors.Wrap(err, "getting limit from rubToCurrentCurrency")
	}

	expenses, err := currency.RubToCurrentCurrency(userCurrency, expRUB, rates)
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
