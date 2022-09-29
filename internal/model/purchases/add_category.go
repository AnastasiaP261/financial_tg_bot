package purchases

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/normalize"
)

func (m *Model) AddCategory(userID int64, category string) error {
	if err := m.Repo.AddCategory(CategoryRow{
		UserID:   userID,
		Category: normalize.Category(category),
	}); err != nil {
		return errors.Wrap(err, "Repo.AddCategory")
	}
	return nil
}
