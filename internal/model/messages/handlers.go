package messages

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func (m *Model) msgReport(msg Message) error {
	res := report.FindStringSubmatch(msg.Text)
	if len(res) < 2 {
		return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
	}

	period, err := m.purchasesModel.ToPeriod(res[1])
	if err != nil {
		return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
	}

	reportTxt, img, err := m.purchasesModel.Report(period, msg.UserID)
	if err != nil {
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}

	if err = m.tgClient.SendMessage(reportTxt, msg.UserID, msg.UserName); err != nil {
		return err
	}

	return m.tgClient.SendImage(img, msg.ChatID, msg.UserName)
}

func (m *Model) msgAddCategory(msg Message) error {
	res := addCategory.FindStringSubmatch(msg.Text)
	if len(res) < 2 {
		return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
	}

	err := m.purchasesModel.AddCategory(msg.UserID, res[1])
	if err != nil {
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}
	return m.tgClient.SendMessage(ScsTxtCategoryAdded, msg.UserID, msg.UserName)
}

func (m *Model) msgAddPurchase(msg Message, sum, category, date string) error {
	if err := m.purchasesModel.AddPurchase(msg.UserID, sum, category, date); err != nil {
		if errors.Is(err, purchases.ErrCategoryNotExist) {
			return m.tgClient.SendMessage(ErrTxtCategoryDoesntExist, msg.UserID, msg.UserName)
		}
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}
	return m.tgClient.SendMessage(ScsTxtPurchaseAdded, msg.UserID, msg.UserName)
}

func (m *Model) msgCurrency(msg Message, rawCY string) error {
	cy, err := m.purchasesModel.StrToCurrency(rawCY)
	if err != nil {
		return m.tgClient.SendMessage(ErrTxtInvalidCurrency, msg.UserID, msg.UserName)
	}

	if err = m.purchasesModel.ChangeUserCurrency(msg.UserID, cy); err != nil {
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}
	return m.tgClient.SendMessage(ScsTxtCurrencyChanged, msg.UserID, msg.UserName)
}
