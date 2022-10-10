package purchases

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrSummaParsing     = errors.New("summa parsing error")
	ErrDateParsing      = errors.New("date parsing error")
	ErrCategoryNotExist = errors.New("such category doesn't exist")
	ErrUnknownPeriod    = errors.New("unknown period")
)

// Repo репозиторий
type Repo interface {
	AddPurchase(req AddPurchaseReq) error
	CategoryExist(req CategoryRow) (bool, error)
	AddCategory(req CategoryRow) error
	GetUserPurchasesFromDate(fromDate time.Time, userID int64) ([]Purchase, error)
	ChangeCurrency(userID int64, currency Currency) error
	GetUserInfo(userID int64) (User, error)
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
