package purchases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	mocks "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases/_mocks"
)

func Test_AddPurchase_OnlySum(t *testing.T) {
	t.Run("целое число", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)

		model := purchases.New(repo, nil, excRateModel)

		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)
		excRateModel.EXPECT().GetExchangeRateToRUB().Return(purchases.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		})
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: purchases.RUB,
		}, nil)

		err := model.AddPurchase(ctx, 123, "123", "", "")
		assert.NoError(t, err)
	})

	t.Run("дробное число", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)

		model := purchases.New(repo, nil, excRateModel)

		excRateModel.EXPECT().GetExchangeRateToRUB().Return(purchases.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		})
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: purchases.RUB,
		}, nil)

		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)

		err := model.AddPurchase(ctx, 123, "234.5", "", "")
		assert.NoError(t, err)
	})

	t.Run("невалидное число", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)

		model := purchases.New(repo, nil, excRateModel)

		err := model.AddPurchase(ctx, 123, "12o.o5", "", "")
		assert.Error(t, err, purchases.ErrSummaParsing)
	})
}

func Test_AddPurchase_SumAndCategory(t *testing.T) {
	t.Run("добавление траты по уже существующей категории", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)

		model := purchases.New(repo, nil, excRateModel)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		excRateModel.EXPECT().GetExchangeRateToRUB().Return(purchases.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		})
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: purchases.RUB,
		}, nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)

		err := model.AddPurchase(ctx, 123, "234.5", "some category", "")
		assert.NoError(t, err)
	})

	t.Run("добавление траты по не существующей категории", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)

		model := purchases.New(repo, nil, excRateModel)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(0), nil)

		err := model.AddPurchase(ctx, 123, "234.5", "some category", "")
		assert.Error(t, err, purchases.ErrCategoryNotExist)
	})
}

func Test_AddPurchase_SumAndCategoryAndDate(t *testing.T) {
	t.Run("добавление с валидной датой", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)

		model := purchases.New(repo, nil, excRateModel)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)
		excRateModel.EXPECT().GetExchangeRateToRUB().Return(purchases.RateToRUB{
			USD: 1,
			EUR: 1,
			CNY: 1,
		})
		repo.EXPECT().GetUserInfo(gomock.Any(), int64(123)).Return(purchases.User{
			UserID:   123,
			Currency: purchases.RUB,
		}, nil)
		repo.EXPECT().AddPurchase(gomock.Any(), gomock.Any()).Return(nil)

		err := model.AddPurchase(ctx, 123, "234.5", "some category", "01.01.2022")
		assert.NoError(t, err)
	})

	t.Run("добавление с не валидной датой", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		repo := mocks.NewMockRepo(ctrl)
		excRateModel := mocks.NewMockExchangeRateGetter(ctrl)

		model := purchases.New(repo, nil, excRateModel)

		repo.EXPECT().GetCategoryID(gomock.Any(), gomock.Any()).Return(uint64(1), nil)

		err := model.AddPurchase(ctx, 123, "234.5", "some category", "01-01-2022")
		assert.Error(t, err, purchases.ErrDateParsing)
	})
}
