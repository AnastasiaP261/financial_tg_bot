package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	"testing"
)

func TestService_ChangeCurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	t.Run("изменение валюты еще не существующего пользователя", func(t *testing.T) {
		err := s.ChangeCurrency(ctx, 123, model.USD)
		assert.NoError(t, err)

		// проверим что запись действительно создалась
		var users []user
		s.db.SelectContext(ctx, &users, "SELECT * FROM users")
		assert.EqualValues(t, []user{{UserID: 123, Currency: USD}}, users)
	})

	t.Run("изменение валюты уже существующего пользователя", func(t *testing.T) {
		err := s.ChangeCurrency(ctx, 123, model.CNY)
		assert.NoError(t, err)

		// проверим что запись действительно создалась
		var users []user
		s.db.SelectContext(ctx, &users, "SELECT * FROM users")
		assert.EqualValues(t, []user{{UserID: 123, Currency: CNY}}, users)
	})
}

func TestService_GetUserInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	// должен сначала создать пользователя, а потом вывести инфу о нем
	userInfo, err := s.GetUserInfo(ctx, 123)
	assert.NoError(t, err)

	// проверим что запись действительно создалась
	var users []user
	s.db.SelectContext(ctx, &users, "SELECT * FROM users")
	assert.EqualValues(t, []user{{UserID: 123, Currency: RUB}}, users)

	assert.Equal(t, model.User{UserID: 123, Currency: model.RUB}, userInfo)
}

func TestService_UserCreateIfNotExist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	err := s.UserCreateIfNotExist(ctx, 123)
	assert.NoError(t, err)

	// проверим что запись действительно создалась
	var users []user
	s.db.SelectContext(ctx, &users, "SELECT * FROM users")
	assert.EqualValues(t, []user{{UserID: 123, Currency: RUB}}, users)
}

func TestService_addUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	err := s.addUser(ctx, 123)
	assert.NoError(t, err)

	// проверим что запись действительно создалась
	var users []user
	s.db.SelectContext(ctx, &users, "SELECT * FROM users")
	assert.EqualValues(t, []user{{UserID: 123, Currency: RUB}}, users)
}

func TestService_getUserInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	// заполнение необходимыми для теста данными
	s.db.ExecContext(ctx, "INSERT INTO users (id, curr) VALUES (123, 'USD')")

	info, err := s.getUserInfo(ctx, 123)
	assert.NoError(t, err)
	assert.Equal(t, user{UserID: 123, Currency: USD}, info)
}

func TestService_userExist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	t.Parallel()

	ctx := context.Background()
	s, close := NewTestDB(ctx, t)
	defer close()

	t.Run("пользователь не существует", func(t *testing.T) {
		ok, err := s.userExist(ctx, 123)

		assert.NoError(t, err)
		assert.False(t, ok)
	})

	// заполнение необходимыми для теста данными
	s.db.ExecContext(ctx, "INSERT INTO users (id, curr) VALUES (123, 'RUB')")

	t.Run("пользователь существует", func(t *testing.T) {
		ok, err := s.userExist(ctx, 123)

		assert.NoError(t, err)
		assert.True(t, ok)
	})
}
