package tg

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/messages"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client *tgbotapi.BotAPI
}

func New(tokenGetter TokenGetter) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SendMessage(text string, userID int64, userName string) error {
	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	log.Printf("[%s] %s", userName, text)

	return nil
}

func (c *Client) SendImage(img []byte, chatID int64, userName string) error {
	b := tgbotapi.FileBytes{Bytes: img}

	_, err := c.client.Send(tgbotapi.NewPhoto(chatID, b))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}

	log.Printf("[%s] image sended", userName)

	return nil
}

func (c *Client) ListenUpdates(ctx context.Context, msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	log.Println("listening for messages")

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			err := msgModel.IncomingMessage(ctx, messages.Message{
				Text:     update.Message.Text,
				UserID:   update.Message.From.ID,
				ChatID:   update.Message.Chat.ID,
				UserName: update.Message.From.UserName,
			})
			if err != nil {
				log.Println("error processing message:", err)
			}
		}
	}
}
