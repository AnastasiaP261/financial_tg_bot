package messages

import (
	"github.com/pkg/errors"
)

func (m *Model) SendMessage(text string, userID int64) error {
	err := m.tgClient.SendMessage(text, userID)
	if err != nil {
		return errors.Wrap(err, "client.SendMessage")
	}

	return nil
}

func (m *Model) SendImage(img []byte, userID int64) error {
	err := m.tgClient.SendImage(img, userID)
	if err != nil {
		return errors.Wrap(err, "client.SendImage")
	}

	return nil
}

func (m *Model) SendKeyboard(text string, userID int64, buttonTexts []string) error {
	err := m.tgClient.SendKeyboard(text, userID, buttonTexts)
	if err != nil {
		return errors.Wrap(err, "client.SendKeyboard")
	}

	return nil
}
