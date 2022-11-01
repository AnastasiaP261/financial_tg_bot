package purchases

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrSummaParsing      = errors.New("summa parsing error")
	ErrLimitParsing      = errors.New("limit parsing error")
	ErrDateParsing       = errors.New("date parsing error")
	ErrCategoryNotExist  = errors.New("such category doesn't exist")
	ErrUnknownPeriod     = errors.New("unknown period")
	ErrInvalidDate       = errors.New("invalid date")
	ErrUserHasntCategory = errors.New("this user hasn't such category")
)

// Repo репозиторий
type Repo interface {
	GetRate(ctx context.Context, y int, m int, d int) (bool, RateToRUB, error)
	AddRate(ctx context.Context, y int, m int, d int, rates RateToRUB) error

	UserCreateIfNotExist(ctx context.Context, userID int64) error
	ChangeCurrency(ctx context.Context, userID int64, currency Currency) error
	GetUserInfo(ctx context.Context, userID int64) (User, error)
	ChangeUserLimit(ctx context.Context, userID int64, newLimit float64) error
	AddCategoryToUser(ctx context.Context, userID int64, catName string) error
	UserHasCategory(ctx context.Context, userID int64, categoryID uint64) (bool, error)

	AddPurchase(ctx context.Context, req AddPurchaseReq) error
	GetUserPurchasesFromDate(ctx context.Context, fromDate time.Time, userID int64) ([]Purchase, error)
	GetUserPurchasesSumFromMonth(ctx context.Context, userID int64, fromDate time.Time) (float64, error)

	GetCategoryID(ctx context.Context, categoryName string) (uint64, error)
	AddCategory(ctx context.Context, categoryName string) error
	GetAllCategories(ctx context.Context) ([]CategoryRow, error)
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
