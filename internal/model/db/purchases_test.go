//go:build test_all || integration_test

package db

import (
	"context"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func Test_AddPurchase(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	fixtures, err := testfixtures.New(
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Files(
			"./../../../test_data/fixtures/users.yml",
			"./../../../test_data/fixtures/categories.yml",
		),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())

	nowTime := time.Now()
	err = s.AddPurchase(ctx, model.AddPurchaseReq{
		UserID:     123,
		Sum:        100,
		CategoryID: 1,
		Date:       nowTime,
		USDRatio:   1,
		CNYRatio:   1,
		EURRatio:   1,
	})

	assert.NoError(t, err)

	// проверим что трата действительно создалась
	var purchases []purchaseTestRow
	selectAllFromTestTablePurchases(ctx, s, &purchases)

	assert.EqualValues(t, []purchaseTestRow{{Sum: 100, CategoryID: 1, USDRatio: 1, CNYRatio: 1, EURRatio: 1}}, purchases)
}

func Test_GetUserPurchasesFromDate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	fixtures, err := testfixtures.New(
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.FilesMultiTables(
			"./../../../test_data/fixtures/get_user_purchases_from_date.yml",
		),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())

	fromTime, _ := time.Parse("02.01.2006", "01.10.2022")
	res, err := s.GetUserPurchasesFromDate(ctx, fromTime, 123)

	assert.NoError(t, err)
	assert.EqualValues(t, []model.Purchase{
		{PurchaseCategory: "some category 1", Summa: 200, RateToRUB: model.RateToRUB{USD: 0.5, EUR: 0.5, CNY: 0.5}},
		{PurchaseCategory: "", Summa: 300, RateToRUB: model.RateToRUB{USD: 0.5, EUR: 0.5, CNY: 0.5}},
		{PurchaseCategory: "some category 2", Summa: 400, RateToRUB: model.RateToRUB{USD: 0.5, EUR: 0.5, CNY: 0.5}},
	}, res)
}

func Test_GetUserPurchasesSumFromMonth(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	fixtures, err := testfixtures.New(
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.FilesMultiTables(
			"./../../../test_data/fixtures/get_user_purchases_sum_from_month.yml",
		),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())

	t.Run("получить сумму трат за 11-ый месяц (до 15 числа включительно)", func(t *testing.T) {
		date, _ := time.Parse("02.01.2006", "15.11.2022")
		res, err := s.GetUserPurchasesSumFromMonth(ctx, 123, date)

		assert.NoError(t, err)
		assert.Equal(t, float64(300), res)
	})

	t.Run("получить сумму трат за 12-ый месяц (до 15 числа включительно)", func(t *testing.T) {
		date, _ := time.Parse("02.01.2006", "15.12.2022")
		res, err := s.GetUserPurchasesSumFromMonth(ctx, 123, date)

		assert.NoError(t, err)
		assert.Equal(t, float64(500), res)
	})
}
