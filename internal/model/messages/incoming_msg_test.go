//go:build test_all || unit_test

package messages

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctx := context.Background()

	sender, purchasesModel, _ := mocksUp(t)
	model := New(sender, purchasesModel, nil)

	sender.EXPECT().SendMessage("hello", int64(123), "name")

	err := model.IncomingMessage(ctx, Message{
		Text:     "/start",
		UserID:   123,
		UserName: "name",
	})

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctx := context.Background()

	sender, purchasesModel, _ := mocksUp(t)
	model := New(sender, purchasesModel, nil)

	sender.EXPECT().SendMessage("Не знаю эту команду", int64(123), "name")

	err := model.IncomingMessage(ctx, Message{
		Text:     "some text",
		UserID:   123,
		UserName: "name",
	})

	assert.NoError(t, err)
}

func Test_OnAddPurchaseCommand(t *testing.T) {
	t.Run("записать только сумму без категории", func(t *testing.T) {
		ctx := context.Background()

		sender, purchasesModel, _ := mocksUp(t)
		model := New(sender, purchasesModel, nil)

		sender.EXPECT().SendMessage("Трата добавлена", int64(123), "name")
		purchasesModel.EXPECT().AddPurchase(gomock.Any(), gomock.Any(), gomock.Any(), "", "").Return(nil)

		err := model.IncomingMessage(ctx, Message{
			Text:     "/add 123.45",
			UserID:   123,
			UserName: "name",
		})

		assert.NoError(t, err)
	})

	t.Run("записать сумму и категорию", func(t *testing.T) {
		ctx := context.Background()

		sender, purchasesModel, _ := mocksUp(t)
		model := New(sender, purchasesModel, nil)

		sender.EXPECT().SendMessage("Трата добавлена", int64(123), "name")
		purchasesModel.EXPECT().AddPurchase(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), "").Return(nil)

		err := model.IncomingMessage(ctx, Message{
			Text:     "/add 123.45 категория какая то",
			UserID:   123,
			UserName: "name",
		})

		assert.NoError(t, err)
	})

	t.Run("записать сумму, категорию и указать дату", func(t *testing.T) {
		ctx := context.Background()

		sender, purchasesModel, _ := mocksUp(t)
		model := New(sender, purchasesModel, nil)

		sender.EXPECT().SendMessage("Трата добавлена", int64(123), "name")
		purchasesModel.EXPECT().AddPurchase(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := model.IncomingMessage(ctx, Message{
			Text:     "/add 123.45 категория какая то 01.01.2022",
			UserID:   123,
			UserName: "name",
		})

		assert.NoError(t, err)
	})
}

func Test_OnAddCategoryCommand(t *testing.T) {
	ctx := context.Background()

	sender, purchasesModel, _ := mocksUp(t)
	model := New(sender, purchasesModel, nil)

	sender.EXPECT().SendMessage("Категория создана", int64(123), "name")
	purchasesModel.EXPECT().AddCategory(gomock.Any(), gomock.Any()).Return(nil)

	err := model.IncomingMessage(ctx, Message{
		Text:     "/category категория",
		UserID:   123,
		UserName: "name",
	})

	assert.NoError(t, err)
}
