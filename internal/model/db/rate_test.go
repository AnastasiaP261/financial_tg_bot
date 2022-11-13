//go:build test_all || integration_test

package db

import (
	"context"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
)

func Test_AddRate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	err := s.AddRate(ctx, 2022, 01, 01, currency.RateToRUB{USD: 2, EUR: 3, CNY: 4})

	assert.NoError(t, err)

	// проверим что курс действительно записался
	var rates []rateTestRow
	selectAllFromTestTableRate(ctx, s, &rates)

	assert.EqualValues(t, []rateTestRow{{Date: 20220101, USD: 2, EUR: 3, CNY: 4}}, rates)
}

func Test_GetRate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	t.Run("курса на эту дату нет", func(t *testing.T) {
		ok, res, err := s.GetRate(ctx, 2022, 01, 01)

		assert.NoError(t, err)
		assert.False(t, ok)
		assert.EqualValues(t, currency.RateToRUB{}, res)
	})

	fixtures, err := testfixtures.New(
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Files(
			"./../../../test_data/fixtures/rate.yml",
		),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())

	t.Run("курс на эту дату есть в базе", func(t *testing.T) {
		ok, res, err := s.GetRate(ctx, 2022, 01, 01)

		assert.NoError(t, err)
		assert.True(t, ok)
		assert.EqualValues(t, currency.RateToRUB{EUR: 2, USD: 3, CNY: 4}, res)
	})
}
