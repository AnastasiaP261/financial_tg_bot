package purchases

import "github.com/pkg/errors"

type User struct {
	UserID   int64
	Currency Currency // выбранная пользователем валюта
}

func (m *Model) ChangeUserCurrency(userID int64, currency Currency) error {
	err := m.Repo.ChangeCurrency(userID, currency)
	if err != nil {
		return errors.Wrap(err, "Repo.ChangeCurrency")
	}
	return nil
}
