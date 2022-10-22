package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

type user struct {
	UserID   uint64   `db:"id"`
	Currency Currency `db:"curr"` // выбранная пользователем валюта
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
func currencyFromModelTypeConv(c model.Currency) (Currency, error) {
	switch c {
	case model.RUB:
		return RUB, nil
	case model.USD:
		return USD, nil
	case model.EUR:
		return EUR, nil
	case model.CNY:
		return CNY, nil
	default:
		return "", errors.New("invalid currency")
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
	res, err := s.getUserInfo(ctx, userID)
	if errors.Is(err, ErrUserDoesntExists) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "getUserInfo")
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

	if _, err = s.db.ExecContext(ctx, q, args...); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

// ChangeCurrency смена валюты пользователя
func (s *Service) ChangeCurrency(ctx context.Context, userID int64, currency model.Currency) error {
	if err := s.UserCreateIfNotExist(ctx, userID); err != nil {
		return errors.Wrap(err, "UserCreateIfNotExist")
	}

	curr, err := currencyFromModelTypeConv(currency)
	if err != nil {
		return errors.Wrap(err, "currencyFromModelTypeConv")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(tblUsers).
		Set(tblUsersColCurrency, curr).
		Where(sq.Eq{tblUsersColID: userID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args...); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

// GetUserInfo возвращает информацию о пользователе в формате модели
func (s *Service) GetUserInfo(ctx context.Context, userID int64) (model.User, error) {
	res, err := s.getUserInfo(ctx, userID)
	if err != nil && errors.Is(err, ErrUserDoesntExists) {

		if err := s.addUser(ctx, userID); err != nil {
			return model.User{}, errors.Wrap(err, "addUser")
		}

		res, err = s.getUserInfo(ctx, userID)
		if err != nil {
			return model.User{}, errors.Wrap(err, "getUserInfo")
		}

	} else if err != nil {
		return model.User{}, errors.Wrap(err, "getUserInfo")
	}

	curr, err := currencyToModelTypeConv(res.Currency)
	if err != nil {
		return model.User{}, errors.Wrap(err, "currencyToModelTypeConv")
	}

	return model.User{
		UserID:   res.UserID,
		Currency: curr,
	}, nil
}

// getUserInfo возвращает информацию о пользователе (для использования внутри пакета)
func (s *Service) getUserInfo(ctx context.Context, userID int64) (user, error) {
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

	rows, err := s.db.QueryxContext(ctx, q, args...)
	if err != nil {
		return user{}, errors.Wrap(err, "db.QueryxContext")
	}

	data := user{}
	if err = readX(rows, &data); err != nil {
		return user{}, errors.Wrap(err, "readX")
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) || data.UserID == 0 {
		return user{}, ErrUserDoesntExists
	} else if err != nil {
		return user{}, errors.Wrap(err, "db.QueryxContext")
	}

	return data, nil
}
