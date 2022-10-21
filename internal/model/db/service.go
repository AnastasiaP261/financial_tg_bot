package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	tblUsers            = "users"
	tblUsersColID       = "id"
	tblUsersColCurrency = "curr"

	tblCategories                = "categories"
	tblCategoriesColID           = "id"
	tblCategoriesColUserID       = "user_id"
	tblCategoriesColCategoryName = "category_name"

	tblPurchases              = "purchases"
	tblPurchasesColID         = "id"
	tblPurchasesColCategoryID = "category_id"
	tblPurchasesColSum        = "sum"
	tblPurchasesColTimestamp  = "ts"
	tblPurchasesColEURRatio   = "eur_ratio"
	tblPurchasesColUSDRatio   = "usd_ratio"
	tblPurchasesColCNYRatio   = "cny_ratio"
)

type configGetter interface {
	DBUri() string
}

type Service struct {
	db *sql.DB
}

func New(ctx context.Context, config configGetter) (*Service, error) {
	db, err := sql.Open("postgres", config.DBUri())
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "db.PingContext")
	}

	return &Service{db: db}, nil
}
