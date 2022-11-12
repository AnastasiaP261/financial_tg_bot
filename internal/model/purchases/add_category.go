package purchases

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/normalize"
)

// AddCategory создать новую категорию
func (m *Model) AddCategory(ctx context.Context, category string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "add category")
	defer span.Finish()

	if err := m.Repo.AddCategory(ctx, normalize.Category(category)); err != nil {
		return errors.Wrap(err, "Repo.AddCategory")
	}
	return nil
}

// GetAllCategories получить все категории
func (m *Model) GetAllCategories(ctx context.Context) ([]CategoryRow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get all categories")
	defer span.Finish()

	res, err := m.Repo.GetAllCategories(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Repo.AddCategory")
	}
	return res, nil
}
