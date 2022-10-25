package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	"testing"
)

func Test_AddRate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	// rate курс валют к RUB
	type rateRow struct {
		Date int64   `json:"date"`
		USD  float64 `db:"usd_ratio"`
		EUR  float64 `db:"eur_ratio"`
		CNY  float64 `db:"cny_ratio"`
	}

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	err := s.AddRate(ctx, 2022, 01, 01, model.RateToRUB{USD: 2, EUR: 3, CNY: 4})

	assert.NoError(t, err)

	// проверим что трата действительно создалась
	var rates []rateRow
	s.db.SelectContext(ctx, &rates, "SELECT *  FROM rate") // nolint:errcheck
	assert.EqualValues(t, []rateRow{{Date: 20220101, USD: 2, EUR: 3, CNY: 4}}, rates)
}

func Test_GetRate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	t.Run("курса на эту дату нет", func(t *testing.T) {
		ok, res, err := s.GetRate(ctx, 2022, 01, 01)

		assert.NoError(t, err)
		assert.False(t, ok)
		assert.EqualValues(t, model.RateToRUB{}, res)
	})

	// заполнение необходимыми для теста данными
	s.db.ExecContext(ctx, "INSERT INTO rate (\"date\", usd_ratio, eur_ratio, cny_ratio) VALUES (20220101, 2, 3, 4)") // nolint:errcheck

	t.Run("курс на эту дату есть в базе", func(t *testing.T) {
		ok, res, err := s.GetRate(ctx, 2022, 01, 01)

		assert.NoError(t, err)
		assert.True(t, ok)
		assert.EqualValues(t, model.RateToRUB{USD: 2, EUR: 3, CNY: 4}, res)
	})
}
