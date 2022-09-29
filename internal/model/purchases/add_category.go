package purchases

import "github.com/pkg/errors"

func (m *Model) AddCategory(userID int64, category string) error {
	if err := m.Repo.AddCategory(CategoryRow{
		UserID:   userID,
		Category: category,
	}); err != nil {
		return errors.Wrap(err, "Repo.AddCategory")
	}
	return nil
}
