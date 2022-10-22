package purchases

import (
	"context"

	"github.com/pkg/errors"
)

type User struct {
	UserID   uint64
	Currency Currency // выбранная пользователем валюта
}

func (m *Model) ChangeUserCurrency(ctx context.Context, userID int64, currency Currency) error {
	err := m.Repo.ChangeCurrency(ctx, userID, currency)
	if err != nil {
		return errors.Wrap(err, "Repo.ChangeCurrency")
	}
	return nil
}
