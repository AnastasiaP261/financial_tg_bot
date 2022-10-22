package db

import (
	"context"

	"github.com/jmoiron/sqlx"
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

	tblPurchases = "purchases"
	//tblPurchasesColID         = "id"
	tblPurchasesColCategoryID = "category_id"
	tblPurchasesColSum        = "sum"
	tblPurchasesColTimestamp  = "ts"
	tblPurchasesColEURRatio   = "eur_ratio"
	tblPurchasesColUSDRatio   = "usd_ratio"
	tblPurchasesColCNYRatio   = "cny_ratio"
)

var (
	ErrUserDoesntExists      = errors.New("user doesnt exists")
	ErrCategoryAlreadyExists = errors.New("category is already exists")
)

type configGetter interface {
	DBUri() string
}

type Service struct {
	db *sqlx.DB
}

func New(ctx context.Context, config configGetter) (*Service, error) {
	db, err := sqlx.Open("postgres", config.DBUri())
	if err != nil {
		return nil, errors.Wrap(err, "sql.Open")
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "db.PingContext")
	}

	return &Service{db: db}, nil
}

type iterableScannerX interface {
	Next() bool
	StructScan(dest any) error
}

// readX позволяет считать в структуру данные полученные из QueryxContext
func readX[T iterableScannerX](rows T, dest any) error {
	for rows.Next() {
		err := rows.StructScan(dest)
		if err != nil {
			return errors.Wrap(err, "rows.StructScan")
		}
	}
	return nil
}

type iterableScanner interface {
	Next() bool
	Scan(dest ...any) error
}

// read позволяет считать в переменные данные полученные из QueryContext
func read[T iterableScanner](rows T, dest any) error {
	for rows.Next() {
		err := rows.Scan(dest)
		if err != nil {
			return errors.Wrap(err, "rows.StructScan")
		}
	}
	return nil
}
