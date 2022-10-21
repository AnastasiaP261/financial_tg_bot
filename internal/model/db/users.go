package db

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

type user struct {
	UserID   uint64
	Currency Currency // выбранная пользователем валюта
}

// Currency тип валюты
type Currency string

const (
	// RUB валюта - рубль
	RUB Currency = "RUB"

	// USD валюта - доллар
	USD Currency = "USD"

	// EUR валюта - евро
	EUR Currency = "EUR"

	// CNY валюта - китайский юань
	CNY Currency = "CNY"
)

func currencyToModelTypeConv(c Currency) (model.Currency, error) {
	switch c {
	case RUB:
		return model.RUB, nil
	case USD:
		return model.USD, nil
	case EUR:
		return model.EUR, nil
	case CNY:
		return model.CNY, nil
	default:
		return 0, errors.New("invalid currency")
	}
}

// UserCreateIfNotExist проверяет, что такой юзер есть в базе, и, если его нет, создает такого юзера.
// Нужно вызывать эту функцию в начале каждой другой команды. Это позволит лениво создать запись о пользователе и
// снимет с модели ответственность за нормализацию данных
func (s *Service) UserCreateIfNotExist(ctx context.Context, userID int64) error {
	ok, err := s.userExist(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "userExist")
	}

	if !ok {
		if err = s.addUser(ctx, userID); err != nil {
			return errors.Wrap(err, "addUser")
		}
	}

	return nil
}

// userExist проверка, что такой юзер уже создан в базе
func (s *Service) userExist(ctx context.Context, userID int64) (bool, error) {
	res, err := s.GetUserInfo(ctx, userID)
	if err != nil {
		return false, errors.Wrap(err, "GetUserInfo")
	}

	if res.UserID == 0 {
		return false, nil
	}

	return true, nil
}

// addUser добавляет юзера с такой айдишкой в базу
func (s *Service) addUser(ctx context.Context, userID int64) error {
	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tblUsers).
		Columns(tblUsersColID, tblUsersColCurrency).
		Values(userID, RUB).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

// ChangeCurrency смена валюты пользователя
func (s *Service) ChangeCurrency(ctx context.Context, userID int64, currency model.Currency) error {
	if err := s.UserCreateIfNotExist(ctx, userID); err != nil {
		return errors.Wrap(err, "UserCreateIfNotExist")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tblUsers).
		Set(tblUsersColCurrency, currency).
		Where(sq.Eq{tblUsersColID: userID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

// GetUserInfo возвращает информацию о пользователе в формате модели
func (s *Service) GetUserInfo(ctx context.Context, userID int64) (model.User, error) {
	res, err := s.getUserInfo(ctx, userID)
	if err != nil {
		return model.User{}, errors.Wrap(err, "getUserInfo")
	}

	curr, err := currencyToModelTypeConv(res.Currency)

	return model.User{
		UserID:   res.UserID,
		Currency: curr,
	}, nil
}

// возвращает информацию о пользователе (для использования внутри пакета)
func (s *Service) getUserInfo(ctx context.Context, userID int64) (user, error) {
	err := s.UserCreateIfNotExist(ctx, userID)
	if err != nil {
		return user{}, errors.Wrap(err, "UserCreateIfNotExist")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tblUsersColID, tblUsersColCurrency).
		From(tblUsers).
		Where(sq.Eq{
			tblUsersColID: userID,
		}).
		ToSql()
	if err != nil {
		return user{}, errors.Wrap(err, "query creating error")
	}

	var (
		id       uint64
		currency Currency
	)

	fmt.Println("### ping", s.db.Ping())

	if err = s.db.QueryRowContext(ctx, q, args).Scan(&id, &currency); err != nil {
		return user{}, errors.Wrap(err, "db.QueryRowContext")
	}

	return user{
		UserID:   id,
		Currency: currency,
	}, nil
}
