package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	"testing"
	"time"
)

func TestService_AddPurchase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	start := func() (context.Context, *Service, func()) {
		ctx := context.Background()
		s, closeFunc := NewTestDB(ctx, t)

		// заполнение необходимыми для теста данными
		s.db.ExecContext(ctx, "INSERT INTO users (id, curr) VALUES (123, 'RUB')")
		s.db.ExecContext(ctx, "INSERT INTO categories (user_id, category_name) VALUES (123, 'some category')")

		return ctx, s, closeFunc
	}

	type purchase struct {
		Sum        float64 `db:"sum"` // сумма траты в рублях
		CategoryID uint64  `db:"category_id"`

		// коэффициенты валют на момент совершения траты
		USDRatio float64 `db:"usd_ratio"`
		CNYRatio float64 `db:"cny_ratio"`
		EURRatio float64 `db:"eur_ratio"`
	}

	t.Run("создание траты с категорией", func(t *testing.T) {
		t.Parallel()
		ctx, s, close := start()
		defer close()

		nowTime := time.Now()
		err := s.AddPurchase(ctx, model.AddPurchaseReq{
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
		var purchases []purchase
		s.db.SelectContext(ctx, &purchases, "SELECT sum, category_id, usd_ratio, cny_ratio, eur_ratio  FROM purchases")
		assert.EqualValues(t, []purchase{{Sum: 100, CategoryID: 1, USDRatio: 1, CNYRatio: 1, EURRatio: 1}}, purchases)
	})

	t.Run("создание траты без категории", func(t *testing.T) {
		t.Parallel()
		ctx, s, close := start()
		defer close()

		nowTime := time.Now()
		err := s.AddPurchase(ctx, model.AddPurchaseReq{
			UserID:     123,
			Sum:        100,
			CategoryID: 0,
			Date:       nowTime,
			USDRatio:   1,
			CNYRatio:   1,
			EURRatio:   1,
		})

		assert.NoError(t, err)

		// проверим что трата действительно создалась
		var purchases []purchase
		s.db.SelectContext(ctx, &purchases, "SELECT sum, category_id, usd_ratio, cny_ratio, eur_ratio  FROM purchases")
		assert.EqualValues(t, []purchase{{Sum: 100, CategoryID: 0, USDRatio: 1, CNYRatio: 1, EURRatio: 1}}, purchases)
	})
}

func TestService_GetUserPurchasesFromDate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	// заполнение необходимыми для теста данными
	s.db.ExecContext(ctx, "INSERT INTO users (id, curr) VALUES (123, 'RUB')")
	s.db.ExecContext(ctx, "INSERT INTO categories (user_id, category_name) VALUES (123, 'some category 1'), (123, 'some category 2')")
	s.db.ExecContext(ctx, "INSERT INTO purchases (category_id, sum, ts, eur_ratio, usd_ratio, cny_ratio) VALUES (1, 100, '2022-09-27', 0.5, 0.5, 0.5), (1, 200,'2022-10-01', 0.5, 0.5, 0.5), (2, 300, '2022-10-24', 0.5, 0.5, 0.5)")

	fromTime, _ := time.Parse("02.01.2006", "01.10.2022")
	res, err := s.GetUserPurchasesFromDate(ctx, fromTime, 123)

	assert.NoError(t, err)
	assert.EqualValues(t, []model.Purchase{
		{PurchaseCategory: "some category 1", Summa: 200, RateToRUB: model.RateToRUB{USD: 0.5, EUR: 0.5, CNY: 0.5}},
		{PurchaseCategory: "some category 2", Summa: 300, RateToRUB: model.RateToRUB{USD: 0.5, EUR: 0.5, CNY: 0.5}},
	}, res)
}
