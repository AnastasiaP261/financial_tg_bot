// nolint: unused
package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/db/_testdb"
)

func newTestDB(ctx context.Context, t *testing.T) (testServ *Service, close func()) {
	container := _testdb.NewTestDatabase(ctx, t)

	time.Sleep(10 * time.Second)

	db, err := New(ctx, container)
	assert.NoError(t, err, "newTestDB")

	return db, func() {
		container.Close(t)
	}
}

func selectAllFromTestTableCategories(ctx context.Context, s *Service, categories *[]category) { // nolint: unused
	_ = s.db.SelectContext(ctx, categories, "SELECT * FROM categories") // nolint:errcheck
}

type purchaseTestRow struct {
	Sum        float64 `db:"sum"` // сумма траты в рублях
	CategoryID uint64  `db:"category_id"`
	UserID     int64   `db:"user_id"`

	// коэффициенты валют на момент совершения траты
	USDRatio float64 `db:"usd_ratio"`
	CNYRatio float64 `db:"cny_ratio"`
	EURRatio float64 `db:"eur_ratio"`
}

func selectAllFromTestTablePurchases(ctx context.Context, s *Service, purchases *[]purchaseTestRow) {
	_ = s.db.SelectContext(ctx, purchases, "SELECT sum, user_id, category_id, usd_ratio, cny_ratio, eur_ratio  FROM purchases") // nolint:errcheck
}

// rate курс валют к RUB
type rateTestRow struct {
	Date int64   `json:"date"`
	USD  float64 `db:"usd_ratio"`
	EUR  float64 `db:"eur_ratio"`
	CNY  float64 `db:"cny_ratio"`
}

func selectAllFromTestTableRate(ctx context.Context, s *Service, rate *[]rateTestRow) {
	_ = s.db.SelectContext(ctx, rate, "SELECT *  FROM rate") // nolint:errcheck
}

func selectAllFromTestTableUsers(ctx context.Context, s *Service, users *[]user) {
	err := s.db.SelectContext(ctx, users, "SELECT * FROM users") // nolint:errcheck
	fmt.Println("### err", err)
}
