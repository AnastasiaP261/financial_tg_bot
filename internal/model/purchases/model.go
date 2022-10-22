package purchases

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrSummaParsing     = errors.New("summa parsing error")
	ErrDateParsing      = errors.New("date parsing error")
	ErrCategoryNotExist = errors.New("such category doesn't exist")
	ErrUnknownPeriod    = errors.New("unknown period")
)

// Repo репозиторий
type Repo interface {
	AddPurchase(ctx context.Context, req AddPurchaseReq) error
	GetCategoryID(ctx context.Context, req CategoryRow) (uint64, error)
	AddCategory(ctx context.Context, req CategoryRow) error
	GetUserPurchasesFromDate(ctx context.Context, fromDate time.Time, userID int64) ([]Purchase, error)
	ChangeCurrency(ctx context.Context, userID int64, currency Currency) error
	GetUserInfo(ctx context.Context, userID int64) (User, error)
}

// ChartDrawer рисовальщик
type ChartDrawer interface {
	// PieChart нарисовать круговую диаграмму трат
	PieChart(data []ReportItem) ([]byte, error)
}

type ExchangeRateGetter interface {
	GetExchangeRateToRUB() RateToRUB
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
