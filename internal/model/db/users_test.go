//go:build test_all || integration_test

package db

import (
	"context"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func Test_ChangeCurrency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	t.Run("изменение валюты еще не существующего пользователя", func(t *testing.T) {
		err := s.ChangeCurrency(ctx, 123, currency.USD)
		assert.NoError(t, err)

		// проверим что запись действительно создалась
		var users []user
		selectAllFromTestTableUsers(ctx, s, &users)

		assert.EqualValues(t, []user{{UserID: 123, Currency: USD, Limit: -1, CategoryIDs: pq.Int64Array{1}}}, users)
	})

	t.Run("изменение валюты уже существующего пользователя", func(t *testing.T) {
		err := s.ChangeCurrency(ctx, 123, currency.CNY)
		assert.NoError(t, err)

		// проверим что запись действительно создалась
		var users []user
		selectAllFromTestTableUsers(ctx, s, &users)

		assert.EqualValues(t, []user{{UserID: 123, Currency: CNY, Limit: -1, CategoryIDs: pq.Int64Array{1}}}, users)
	})
}

func Test_GetUserInfo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	// должен сначала создать пользователя, а потом вывести инфу о нем
	userInfo, err := s.GetUserInfo(ctx, 123)
	assert.NoError(t, err)

	// проверим что запись действительно создалась
	var users []user
	selectAllFromTestTableUsers(ctx, s, &users)

	assert.EqualValues(t, []user{{UserID: 123, Currency: RUB, Limit: -1, CategoryIDs: pq.Int64Array{1}}}, users)
	assert.Equal(t, model.User{UserID: 123, Currency: currency.RUB, Limit: -1, Categories: []int64{1}}, userInfo)
}

func Test_UserCreateIfNotExist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	err := s.UserCreateIfNotExist(ctx, 123)
	assert.NoError(t, err)

	// проверим что запись действительно создалась
	var users []user
	selectAllFromTestTableUsers(ctx, s, &users)

	assert.EqualValues(t, []user{{UserID: 123, Currency: RUB, Limit: -1, CategoryIDs: pq.Int64Array{1}}}, users)
}

func Test_addUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	err := s.addUser(ctx, 123)
	assert.NoError(t, err)

	// проверим что запись действительно создалась
	var users []user
	selectAllFromTestTableUsers(ctx, s, &users)

	assert.EqualValues(t, []user{{UserID: 123, Currency: RUB, Limit: -1, CategoryIDs: pq.Int64Array{1}}}, users)
}

func Test_getUserInfo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	fixtures, err := testfixtures.New(
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Files(
			"./../../../test_data/fixtures/users.yml",
		),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())

	info, err := s.getUserInfo(ctx, 123)
	assert.NoError(t, err)
	assert.Equal(t, user{UserID: 123, Currency: RUB, Limit: -1, CategoryIDs: pq.Int64Array{1, 2}}, info)
}

func Test_userExist(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	t.Run("пользователь не существует", func(t *testing.T) {
		ok, err := s.userExist(ctx, 123)

		assert.NoError(t, err)
		assert.False(t, ok)
	})

	fixtures, err := testfixtures.New(
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Files(
			"./../../../test_data/fixtures/users.yml",
		),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())

	t.Run("пользователь существует", func(t *testing.T) {
		ok, err := s.userExist(ctx, 123)

		assert.NoError(t, err)
		assert.True(t, ok)
	})
}

func Test_ChangeUserLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	t.Run("изменение месячного лимита еще не существующего пользователя", func(t *testing.T) {
		err := s.ChangeUserLimit(ctx, 123, 100)
		assert.NoError(t, err)

		// проверим что запись действительно создалась
		var users []user
		selectAllFromTestTableUsers(ctx, s, &users)

		assert.EqualValues(t, []user{{UserID: 123, Currency: RUB, Limit: 100, CategoryIDs: pq.Int64Array{1}}}, users)
	})

	t.Run("изменение месячного лимита уже существующего пользователя", func(t *testing.T) {
		err := s.ChangeUserLimit(ctx, 123, 200)
		assert.NoError(t, err)

		// проверим что запись действительно создалась
		var users []user
		selectAllFromTestTableUsers(ctx, s, &users)

		assert.EqualValues(t, []user{{UserID: 123, Currency: RUB, Limit: 200, CategoryIDs: pq.Int64Array{1}}}, users)
	})
}

func Test_AddCategoryToUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	fixtures, err := testfixtures.New(
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.FilesMultiTables(
			"./../../../test_data/fixtures/add_category_to_user.yml",
		),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())

	err = s.AddCategoryToUser(ctx, 123, "some category")

	assert.NoError(t, err)

	// проверим что запись действительно изменилась
	var users []user
	selectAllFromTestTableUsers(ctx, s, &users)

	assert.EqualValues(t, []user{{UserID: 123, Currency: RUB, Limit: -1, CategoryIDs: pq.Int64Array{1, 2}}}, users)
}

func Test_UserHasCategory(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	fixtures, err := testfixtures.New(
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Files(
			"./../../../test_data/fixtures/users.yml",
			"./../../../test_data/fixtures/categories.yml",
		),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())

	t.Run("у юзера есть категория", func(t *testing.T) {
		has, err := s.UserHasCategory(ctx, 123, 2)

		assert.NoError(t, err)
		assert.True(t, has)
	})
	t.Run("у юзера нет категории", func(t *testing.T) {
		has, err := s.UserHasCategory(ctx, 123, 3)

		assert.NoError(t, err)
		assert.False(t, has)
	})
}

func TestService_GetUserCategories(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s, close := newTestDB(ctx, t)
	defer close()

	fixtures, err := testfixtures.New(
		testfixtures.Database(s.db.DB),
		testfixtures.Dialect("postgres"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Files(
			"./../../../test_data/fixtures/users.yml",
			"./../../../test_data/fixtures/categories.yml",
		),
	)
	assert.NoError(t, err)
	assert.NoError(t, fixtures.Load())

	data, err := s.GetUserCategories(ctx, 123)

	assert.NoError(t, err)
	assert.EqualValues(t, []string{"Не заданная категория", "some category"}, data)
}
