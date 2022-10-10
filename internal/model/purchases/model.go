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
	GetReport(fromDate time.Time, userID int64) ([]ReportItem, error)
}

// ChartDrawer рисовальщик
type ChartDrawer interface {
	// PieChart нарисовать круговую диаграмму трат
	PieChart(data []ReportItem) ([]byte, error)
}

type Model struct {
	Repo        Repo
	ChartDrawer ChartDrawer
}

func New(repo Repo, drawer ChartDrawer) *Model {
	return &Model{Repo: repo, ChartDrawer: drawer}
}
