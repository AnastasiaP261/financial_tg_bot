//go:build test_all || unit_test

package purchases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	mocks "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases/_mocks"
)

func Test_AddPurchase_OnlySum(t *testing.T) {
	t.Run("целое число", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, nil, nil)

		excRateModel.EXPECT().GetExchangeRateToRUB().Return(currency.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		})
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: currency.RUB,
			Limit:    -1,
		}, nil)
		repo.EXPECT().GetUserPurchasesSumFromMonth(gomock.Any(), int64(123), gomock.Any()).Return(float64(100), nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)
		redis.EXPECT().Delete(gomock.Any(), "123report")

		res, err := model.AddPurchase(ctx, 123, "123", "", "")

		assert.NoError(t, err)
		assert.Equal(t, purchases.ExpensesAndLimit{
			Limit:         -1,
			Expenses:      100 + 123,
			Currency:      currency.RUB,
			LimitExceeded: false,
		}, res)
	})

	t.Run("дробное число", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, nil, nil)

		excRateModel.EXPECT().GetExchangeRateToRUB().Return(currency.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		})
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: currency.RUB,
			Limit:    -1,
		}, nil)
		repo.EXPECT().GetUserPurchasesSumFromMonth(gomock.Any(), int64(123), gomock.Any()).Return(float64(100), nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)
		redis.EXPECT().Delete(gomock.Any(), "123report")

		res, err := model.AddPurchase(ctx, 123, "234.5", "", "")

		assert.NoError(t, err)
		assert.Equal(t, purchases.ExpensesAndLimit{
			Limit:         -1,
			Expenses:      100 + 234.5,
			Currency:      currency.RUB,
			LimitExceeded: false,
		}, res)
	})

	t.Run("невалидное число", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, redis, nil)

		_, err := model.AddPurchase(ctx, 123, "12o.o5", "", "")
		assert.Error(t, err, purchases.ErrSummaParsing)
	})
}

func Test_AddPurchase_SumAndCategory(t *testing.T) {
	t.Run("добавление траты по уже существующей категории", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, nil, nil)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		repo.EXPECT().UserHasCategory(gomock.Any(), int64(123), uint64(1)).Return(true, nil)
		excRateModel.EXPECT().GetExchangeRateToRUB().Return(currency.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		})
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: currency.RUB,
			Limit:    -1,
		}, nil)
		repo.EXPECT().GetUserPurchasesSumFromMonth(gomock.Any(), int64(123), gomock.Any()).Return(float64(100), nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)
		redis.EXPECT().Delete(gomock.Any(), "123report")

		res, err := model.AddPurchase(ctx, 123, "234.5", "some category", "")

		assert.NoError(t, err)
		assert.Equal(t, purchases.ExpensesAndLimit{
			Limit:         -1,
			Expenses:      100 + 234.5,
			Currency:      currency.RUB,
			LimitExceeded: false,
		}, res)
	})

	t.Run("добавление траты по не существующей категории", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, redis, nil)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(0), nil)

		_, err := model.AddPurchase(ctx, 123, "234.5", "some category", "")
		assert.Error(t, err, purchases.ErrCategoryNotExist)
	})
}

func Test_AddPurchase_SumAndCategoryAndDate(t *testing.T) {
	t.Run("добавление с валидной датой", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, redis, nil)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		repo.EXPECT().UserHasCategory(gomock.Any(), int64(123), uint64(1)).Return(true, nil)
		repo.EXPECT().GetRate(gomock.Any(), 2022, 1, 1).Return(true, currency.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		}, nil)
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: currency.RUB,
			Limit:    -1,
		}, nil)
		repo.EXPECT().GetUserPurchasesSumFromMonth(gomock.Any(), int64(123), gomock.Any()).Return(float64(100), nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)
		redis.EXPECT().Delete(gomock.Any(), "123report")

		_, err := model.AddPurchase(ctx, 123, "234.5", "some category", "01.01.2022")
		assert.NoError(t, err)
	})

	t.Run("добавление с не валидной датой", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, redis, nil)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		repo.EXPECT().UserHasCategory(gomock.Any(), int64(123), uint64(1)).Return(true, nil)

		_, err := model.AddPurchase(ctx, 123, "234.5", "some category", "01-01-2022")
		assert.Error(t, err, purchases.ErrDateParsing)
	})
}

func Test_AddPurchase_Limits(t *testing.T) {
	t.Run("у юзера не установлен лимит", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, redis, nil)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		repo.EXPECT().UserHasCategory(gomock.Any(), int64(123), uint64(1)).Return(true, nil)
		repo.EXPECT().GetRate(gomock.Any(), 2022, 1, 1).Return(true, currency.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		}, nil)
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: currency.RUB,
			Limit:    -1,
		}, nil)
		repo.EXPECT().GetUserPurchasesSumFromMonth(gomock.Any(), int64(123), gomock.Any()).Return(float64(500), nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)
		redis.EXPECT().Delete(gomock.Any(), "123report")

		expAndLim, err := model.AddPurchase(ctx, 123, "234.5", "some category", "01.01.2022")

		assert.NoError(t, err)
		assert.Equal(t, purchases.ExpensesAndLimit{
			Limit:         -1,
			Expenses:      500 + 234.5,
			Currency:      currency.RUB,
			LimitExceeded: false,
		}, expAndLim)
	})

	t.Run("лимит установлен, не превышен", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, redis, nil)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		repo.EXPECT().UserHasCategory(gomock.Any(), int64(123), uint64(1)).Return(true, nil)
		repo.EXPECT().GetRate(gomock.Any(), 2022, 1, 1).Return(true, currency.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		}, nil)
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: currency.RUB,
			Limit:    1000,
		}, nil)
		repo.EXPECT().GetUserPurchasesSumFromMonth(gomock.Any(), int64(123), gomock.Any()).Return(float64(500), nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)
		redis.EXPECT().Delete(gomock.Any(), "123report")

		expAndLim, err := model.AddPurchase(ctx, 123, "234.5", "some category", "01.01.2022")

		assert.NoError(t, err)
		assert.Equal(t, purchases.ExpensesAndLimit{
			Limit:         1000,
			Expenses:      500 + 234.5,
			Currency:      currency.RUB,
			LimitExceeded: false,
		}, expAndLim)
	})

	t.Run("лимит установлен, превышен", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, redis, nil)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		repo.EXPECT().UserHasCategory(gomock.Any(), int64(123), uint64(1)).Return(true, nil)
		repo.EXPECT().GetRate(gomock.Any(), 2022, 1, 1).Return(true, currency.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		}, nil)
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: currency.RUB,
			Limit:    1000,
		}, nil)
		repo.EXPECT().GetUserPurchasesSumFromMonth(gomock.Any(), int64(123), gomock.Any()).Return(float64(800), nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)
		redis.EXPECT().Delete(gomock.Any(), "123report")

		expAndLim, err := model.AddPurchase(ctx, 123, "234.5", "some category", "01.01.2022")

		assert.NoError(t, err)
		assert.Equal(t, purchases.ExpensesAndLimit{
			Limit:         1000,
			Expenses:      800 + 234.5,
			Currency:      currency.RUB,
			LimitExceeded: true,
		}, expAndLim)
	})

	t.Run("лимит установлен, не превышен, основная валюта не рубль", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)
		redis := mocks.NewMockReportsStore(ctrl)

		model := purchases.New(repo, excRateModel, redis, nil)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		repo.EXPECT().UserHasCategory(gomock.Any(), int64(123), uint64(1)).Return(true, nil)
		repo.EXPECT().GetRate(gomock.Any(), 2022, 1, 1).Return(true, currency.RateToRUB{
			USD: 2,
			EUR: 2,
			CNY: 2,
		}, nil)
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: currency.USD,
			Limit:    1000,
		}, nil)
		repo.EXPECT().GetUserPurchasesSumFromMonth(gomock.Any(), int64(123), gomock.Any()).Return(float64(500), nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)
		redis.EXPECT().Delete(gomock.Any(), "123report")

		expAndLim, err := model.AddPurchase(ctx, 123, "234.5", "some category", "01.01.2022")

		assert.NoError(t, err)
		assert.Equal(t, purchases.ExpensesAndLimit{
			Limit:         1000 * 2,
			Expenses:      500*2 + 234.5,
			Currency:      currency.USD,
			LimitExceeded: false,
		}, expAndLim)
	})
}
