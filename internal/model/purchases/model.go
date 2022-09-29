package purchases

import (
	"github.com/pkg/errors"
	"time"
)

var (
	ErrSummaParsing     = errors.New("summa parsing error")
	ErrDateParsing      = errors.New("date parsing error")
	ErrCategoryNotExist = errors.New("such category doesn't exist")
	ErrUnknownPeriod    = errors.New("unknown period")
)

type Repo interface {
	AddPurchase(req AddPurchaseReq) error
	CategoryExist(req CategoryRow) (bool, error)
	AddCategory(req CategoryRow) error
	GetReport(fromDate time.Time, userID int64) ([]ReportItem, error)
}

type Model struct {
	Repo Repo
}

func New(repo Repo) *Model {
	return &Model{Repo: repo}
}
