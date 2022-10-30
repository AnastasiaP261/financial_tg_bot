//go:build test_all || integration_test

package db

import (
	"context"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func Test_GetCategoryID(t *testing.T) {
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

	t.Run("категория существует", func(t *testing.T) {
		id, err := s.GetCategoryID(ctx, model.CategoryRow{
			UserID:   123,
			Category: "some category",
		})

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), id)
	})

	t.Run("категория не существует", func(t *testing.T) {
		id, err := s.GetCategoryID(ctx, model.CategoryRow{
			UserID:   123,
			Category: "not existing category",
		})

		assert.NoError(t, err)
		assert.Equal(t, uint64(0), id)
	})
}

func Test_AddCategory(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	t.Run("добавление категории", func(t *testing.T) {
		err := s.AddCategory(ctx, model.CategoryRow{
			UserID:   123,
			Category: "some category",
		})

		assert.NoError(t, err)

		// проверим что категория действительно создалась
		var categories []category
		selectAllFromTestTableCategories(ctx, s, &categories)

		assert.EqualValues(t, []category{{1, 123, "some category"}}, categories)
	})

	t.Run("добавление существующей категории", func(t *testing.T) {
		err := s.AddCategory(ctx, model.CategoryRow{
			UserID:   123,
			Category: "some category",
		})

		assert.ErrorIs(t, err, ErrCategoryAlreadyExists)

		// проверим что лишняя категория не создалась
		var categories []category
		selectAllFromTestTableCategories(ctx, s, &categories)

		assert.EqualValues(t, []category{{1, 123, "some category"}}, categories)
	})
}
