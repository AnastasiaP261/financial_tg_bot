package store

import (
	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

// UserCreateIfNotExist проверяет, что такой юзер есть в базе, и, если его нет, создает такого юзера.
// Нужно вызывать эту функцию в начале каждой другой команды. Это позволит лениво создать запись о пользователе и
// снимет с модели ответственность за нормализацию данных
func (s *Service) UserCreateIfNotExist(userID int64) error {
	ok, err := s.userExist(userID)
	if err != nil {
		return errors.Wrap(err, "userExist")
	}

	if !ok {
		if err = s.addUser(userID); err != nil {
			return errors.Wrap(err, "addUser")
		}
	}

	return nil
}

// userExist проверка, что такой юзер уже создан в базе
func (s *Service) userExist(userID int64) (bool, error) {
	for _, row := range s.usersAccessRead() {
		if row.UserID == userID {
			return true, nil
		}
	}

	return false, nil
}

// addUser добавляет юзера с такой айдишкой в базу
func (s *Service) addUser(userID int64) error {
	s.usersAccessWrite(user{
		UserID:   userID,
		Currency: RUB,
	}, flagAdd)
	return nil
}

// ChangeCurrency смена валюты пользователя
func (s *Service) ChangeCurrency(userID int64, currency model.Currency) error {
	if err := s.UserCreateIfNotExist(userID); err != nil {
		return errors.Wrap(err, "UserCreateIfNotExist")
	}

	for _, row := range s.usersAccessRead() {
		if row.UserID == userID {
			s.usersAccessWrite(user{
				UserID:   userID,
				Currency: Currency(currency),
			}, flagEditCurrency)
		}
	}

	return nil
}

// GetUserInfo возвращает информацию о пользователе в формате модели
func (s *Service) GetUserInfo(userID int64) (model.User, error) {
	res, err := s.getUserInfo(userID)
	if err != nil {
		return model.User{}, errors.Wrap(err, "getUserInfo")
	}

	return model.User{
		UserID:   res.UserID,
		Currency: model.Currency(res.Currency),
	}, nil
}

// возвращает информацию о пользователе (для использования внутри пакета)
func (s *Service) getUserInfo(userID int64) (user, error) {
	err := s.UserCreateIfNotExist(userID)
	if err != nil {
		return user{}, errors.Wrap(err, "UserCreateIfNotExist")
	}

	for _, row := range s.usersAccessRead() {
		if row.UserID == userID {
			return row, nil
		}
	}

	return user{}, errors.New("user doesn't exists")
}
