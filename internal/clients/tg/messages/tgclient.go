package messages

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/clients/tg"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/messages"
)

type TokenGetter interface {
	Token() string
}

type MsgHandler struct {
	client *tgbotapi.BotAPI
}

func New(conf TokenGetter) (*MsgHandler, error) {
	client, err := tgbotapi.NewBotAPI(conf.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &MsgHandler{
		client: client,
	}, nil
}

func (m *MsgHandler) SendMessage(text string, userID int64) error {
	_, err := m.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}

	return nil
}

func (m *MsgHandler) SendImage(img []byte, userId int64) error {
	b := tgbotapi.FileBytes{Bytes: img}

	_, err := m.client.Send(tgbotapi.NewPhoto(userId, b))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}

	return nil
}

func (m *MsgHandler) SendKeyboard(text string, userId int64, buttonTexts []string) error {
	msg := tgbotapi.NewMessage(userId, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, txt := range buttonTexts {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(txt, txt)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard

	_, err := m.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}

	return nil
}

func (m *MsgHandler) IncomingCallback(ctx context.Context, model tg.MsgModel, update tgbotapi.Update) error {
	return model.IncomingCallback(ctx, messages.Callback{
		UserID:   update.CallbackQuery.From.ID,
		UserName: update.CallbackQuery.From.UserName,
		Data:     update.CallbackQuery.Data,
	})
}

func (m *MsgHandler) IncomingMessage(ctx context.Context, model tg.MsgModel, update tgbotapi.Update) error {
	return model.IncomingMessage(ctx, messages.Message{
		Text:     update.Message.Text,
		UserID:   update.Message.From.ID,
		UserName: update.Message.From.UserName,
	})
}
