package purchases

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrSummaParsing     = errors.New("summa parsing error")
	ErrLimitParsing     = errors.New("limit parsing error")
	ErrDateParsing      = errors.New("date parsing error")
	ErrCategoryNotExist = errors.New("such category doesn't exist")
	ErrUnknownPeriod    = errors.New("unknown period")
	ErrInvalidDate      = errors.New("invalid date")
)

// Repo репозиторий
type Repo interface {
	AddPurchase(ctx context.Context, req AddPurchaseReq) error
	GetCategoryID(ctx context.Context, req CategoryRow) (uint64, error)
	AddCategory(ctx context.Context, req CategoryRow) error
	GetUserPurchasesFromDate(ctx context.Context, fromDate time.Time, userID int64) ([]Purchase, error)
	ChangeCurrency(ctx context.Context, userID int64, currency Currency) error
	GetUserInfo(ctx context.Context, userID int64) (User, error)
	GetRate(ctx context.Context, y, m, d int) (bool, RateToRUB, error)
	AddRate(ctx context.Context, y, m, d int, rates RateToRUB) error
	ChangeUserLimit(ctx context.Context, userID int64, newLimit float64) error
	GetUserPurchasesSumFromMonth(ctx context.Context, userID int64, date time.Time) (float64, error)
}

// ChartDrawer рисовальщик
type ChartDrawer interface {
	// PieChart нарисовать круговую диаграмму трат
	PieChart(data []ReportItem) ([]byte, error)
}

type ExchangeRateGetter interface {
	GetExchangeRateToRUB() RateToRUB
	GetExchangeRateToRUBFromDate(ctx context.Context, y, m, d int) (RateToRUB, error)
}

type Model struct {
	Repo               Repo
	ChartDrawer        ChartDrawer
	ExchangeRatesModel ExchangeRateGetter
}

func New(repo Repo, drawer ChartDrawer, exchangeRatesModel ExchangeRateGetter) *Model {
	return &Model{
		Repo:               repo,
		ChartDrawer:        drawer,
		ExchangeRatesModel: exchangeRatesModel,
	}
}
