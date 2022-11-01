package purchases

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/normalize"
)

// AddCategory создать новую категорию
func (m *Model) AddCategory(ctx context.Context, category string) error {
	if err := m.Repo.AddCategory(ctx, normalize.Category(category)); err != nil {
		return errors.Wrap(err, "Repo.AddCategory")
	}
	return nil
}

// GetAllCategories получить все категории
func (m *Model) GetAllCategories(ctx context.Context) ([]CategoryRow, error) {
	res, err := m.Repo.GetAllCategories(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Repo.AddCategory")
	}
	return res, nil
}
