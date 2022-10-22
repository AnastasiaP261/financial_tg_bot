package purchases

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/normalize"
)

// AddCategory добавить новую категорию
func (m *Model) AddCategory(ctx context.Context, userID int64, category string) error {
	if err := m.Repo.AddCategory(ctx, CategoryRow{
		UserID:   userID,
		Category: normalize.Category(category),
	}); err != nil {
		return errors.Wrap(err, "Repo.AddCategory")
	}
	return nil
}
