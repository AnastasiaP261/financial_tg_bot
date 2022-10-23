//go:build test_all || integration_test

package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func TestService_GetCategoryID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	// заполнение необходимыми для теста данными
	s.db.ExecContext(ctx, "INSERT INTO users (id, curr) VALUES (123, 'RUB')")                              // nolint:errcheck
	s.db.ExecContext(ctx, "INSERT INTO categories (user_id, category_name) VALUES (123, 'some category')") // nolint:errcheck

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

func TestService_AddCategory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	t.Run("добавление категории", func(t *testing.T) {
		err := s.AddCategory(ctx, model.CategoryRow{
			UserID:   123,
			Category: "some category",
		})

		assert.NoError(t, err)

		// проверим что категория действительно создалась
		var categories []category
		s.db.SelectContext(ctx, &categories, "SELECT * FROM categories") // nolint:errcheck

		assert.EqualValues(t, []category{{1, 123, "some category"}}, categories)
	})

	t.Run("добавление существующей категории", func(t *testing.T) {
		err := s.AddCategory(ctx, model.CategoryRow{
			UserID:   123,
			Category: "some category",
		})
		t.Parallel()

		assert.ErrorIs(t, err, ErrCategoryAlreadyExists)

		// проверим что лишняя категория не создалась
		var categories []category
		s.db.SelectContext(ctx, &categories, "SELECT * FROM categories") // nolint:errcheck

		assert.Nil(t, categories)
	})
}
