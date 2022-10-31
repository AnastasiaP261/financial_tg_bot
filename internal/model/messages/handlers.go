package messages

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func (m *Model) msgReport(ctx context.Context, msg Message) error {
	res := report.FindStringSubmatch(msg.Text)
	if len(res) < 2 {
		return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
	}

	period, err := m.purchasesModel.ToPeriod(res[1])
	if err != nil {
		return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
	}

	reportTxt, img, err := m.purchasesModel.Report(ctx, period, msg.UserID)
	if err != nil {
		err = errors.Wrap(err, "purchasesModel.Report")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}

	if err = m.tgClient.SendMessage(reportTxt, msg.UserID, msg.UserName); err != nil {
		err = errors.Wrap(err, "tgClient.SendMessage")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}

	return m.tgClient.SendImage(img, msg.ChatID, msg.UserName)
}

func (m *Model) msgAddCategory(ctx context.Context, msg Message) error {
	res := addCategory.FindStringSubmatch(msg.Text)
	if len(res) < 2 {
		return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
	}

	err := m.purchasesModel.AddCategory(ctx, msg.UserID, res[1])
	if err != nil {
		err = errors.Wrap(err, "purchasesModel.AddCategory")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}
	return m.tgClient.SendMessage(ScsTxtCategoryAdded, msg.UserID, msg.UserName)
}

func (m *Model) msgAddPurchase(ctx context.Context, msg Message, sum, category, date string) error {
	expAndLim, err := m.purchasesModel.AddPurchase(ctx, msg.UserID, sum, category, date)
	if err != nil {
		err = errors.Wrap(err, "purchasesModel.AddPurchase")
		if errors.Is(err, purchases.ErrCategoryNotExist) {
			return m.tgClient.SendMessage(ErrTxtCategoryDoesntExist, msg.UserID, msg.UserName)
		}
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}

	txt := ScsTxtPurchaseAdded
	if expAndLim.Limit != -1 {
		txt += fmt.Sprintf("\n\nУ вас установлен лимит: %.2f %s. За этот месяц вы потратили уже %.2f %s.",
			expAndLim.Limit, string(expAndLim.Currency), expAndLim.Limit, string(expAndLim.Currency))
		if expAndLim.LimitExceeded {
			txt += "\nВЫ ПРЕВЫСИЛИ ЛИМИТ!"
		}
	}

	return m.tgClient.SendMessage(txt, msg.UserID, msg.UserName)
}

func (m *Model) msgCurrency(ctx context.Context, msg Message, rawCY string) error {
	cy, err := m.purchasesModel.StrToCurrency(rawCY)
	if err != nil {
		return m.tgClient.SendMessage(ErrTxtInvalidCurrency, msg.UserID, msg.UserName)
	}

	if err = m.purchasesModel.ChangeUserCurrency(ctx, msg.UserID, cy); err != nil {
		err = errors.Wrap(err, "purchasesModel.ChangeUserCurrency")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}
	return m.tgClient.SendMessage(ScsTxtCurrencyChanged, msg.UserID, msg.UserName)
}

func (m *Model) msgLimit(ctx context.Context, msg Message, limit string) error {
	if err := m.purchasesModel.ChangeUserLimit(ctx, msg.UserID, limit); err != nil {
		if errors.Is(err, purchases.ErrLimitParsing) {
			return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}
		err = errors.Wrap(err, "purchasesModel.ChangeUserLimit")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}
	return m.tgClient.SendMessage(ScsTxtCurrencyChanged, msg.UserID, msg.UserName)
}
