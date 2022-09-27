package messages

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/messages/_mocks"
)

func mocksUp(t *testing.T) (*mocks.MockMessageSender, *mocks.MockPurchasesModel) {
	ctrl := gomock.NewController(t)
	return mocks.NewMockMessageSender(ctrl), mocks.NewMockPurchasesModel(ctrl)
}

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	sender, purchasesModel := mocksUp(t)
	model := New(sender, purchasesModel)

	sender.EXPECT().SendMessage("hello", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	sender, purchasesModel := mocksUp(t)
	model := New(sender, purchasesModel)

	sender.EXPECT().SendMessage("Не знаю эту команду", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "some text",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAddPuchaseCommand(t *testing.T) {
	t.Run("записать только сумму без категории", func(t *testing.T) {
		sender, purchasesModel := mocksUp(t)
		model := New(sender, purchasesModel)

		sender.EXPECT().SendMessage("Трата добавлена", int64(123))
		purchasesModel.EXPECT().AddPurchase(gomock.Any(), gomock.Any(), "", "").Return(nil)

		err := model.IncomingMessage(Message{
			Text:   "/add 123.45",
			UserID: 123,
		})

		assert.NoError(t, err)
	})

	t.Run("записать сумму и категорию", func(t *testing.T) {
		sender, purchasesModel := mocksUp(t)
		model := New(sender, purchasesModel)

		sender.EXPECT().SendMessage("Трата добавлена", int64(123))
		purchasesModel.EXPECT().AddPurchase(gomock.Any(), gomock.Any(), gomock.Any(), "").Return(nil)

		err := model.IncomingMessage(Message{
			Text:   "/add 123.45 категория какая то",
			UserID: 123,
		})

		assert.NoError(t, err)
	})

	t.Run("записать сумму, категорию и указать дату", func(t *testing.T) {
		sender, purchasesModel := mocksUp(t)
		model := New(sender, purchasesModel)

		sender.EXPECT().SendMessage("Трата добавлена", int64(123))
		purchasesModel.EXPECT().AddPurchase(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := model.IncomingMessage(Message{
			Text:   "/add 123.45 категория какая то 01.01.2022",
			UserID: 123,
		})

		assert.NoError(t, err)
	})
}
