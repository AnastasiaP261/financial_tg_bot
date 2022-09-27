package purchases

import "github.com/pkg/errors"

var (
	ErrSummaParsing         = errors.New("summa parsing error")
	ErrDateParsing          = errors.New("date parsing error")
	ErrCategoryExistRequest = errors.New("category_exist request error")
	ErrCategoryNotExist     = errors.New("such category doesn't exist")
)

type Repo interface {
	AddPurchase(req AddPurchaseReq) error
	CategoryExist(req CategoryExistReq) (bool, error)
}

type Model struct {
	Repo Repo
}

func New(repo Repo) *Model {
	return &Model{Repo: repo}
}
