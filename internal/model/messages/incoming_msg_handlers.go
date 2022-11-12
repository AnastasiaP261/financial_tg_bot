package messages

import (
	"context"
	"fmt"
	"sort"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func (m *Model) msgReport(ctx context.Context, Send Message) error {
	res := report.FindStringSubmatch(Send.Text)
	if len(res) < 2 {
		return m.tgClient.SendMessage(ErrTxtInvalidInput, Send.UserID)
	}

	period, err := m.purchasesModel.ToPeriod(res[1])
	if err != nil {
		return m.tgClient.SendMessage(ErrTxtInvalidInput, Send.UserID)
	}

	reportTxt, img, err := m.purchasesModel.Report(ctx, period, Send.UserID)
	if err != nil {
		err = errors.Wrap(err, "purchasesModel.Report")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), Send.UserID)
	}

	if err = m.tgClient.SendMessage(reportTxt, Send.UserID); err != nil {
		err = errors.Wrap(err, "tgClient.SendMessage")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), Send.UserID)
	}

	return m.tgClient.SendImage(img, Send.UserID)
}

func (m *Model) msgAddCategory(ctx context.Context, Send Message) error {
	res := addCategory.FindStringSubmatch(Send.Text)
	if len(res) < 2 {
		return m.tgClient.SendMessage(ErrTxtInvalidInput, Send.UserID)
	}

	err := m.purchasesModel.AddCategory(ctx, res[1])
	if err != nil {
		err = errors.Wrap(err, "purchasesModel.AddCategory")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), Send.UserID)
	}
	return m.tgClient.SendMessage(ScsTxtCategoryCreated, Send.UserID)
}

func (m *Model) msgAddPurchase(ctx context.Context, Send Message, sum, category, date string) error {
	expAndLim, err := m.purchasesModel.AddPurchase(ctx, Send.UserID, sum, category, date)
	if err != nil {
		err = errors.Wrap(err, "purchasesModel.AddPurchase")

		if errors.Is(err, purchases.ErrCategoryNotExist) || errors.Is(err, purchases.ErrUserHasntCategory) {
			categories, err := m.purchasesModel.GetAllCategories(ctx)
			if err != nil {
				return m.tgClient.SendMessage("Ошибочка: "+err.Error(), Send.UserID)
			}

			buttons := make([]string, len(categories))
			for i := range categories {
				buttons[i] = categories[i].Category
			}
			sort.Strings(buttons)
			buttons = append(buttons, ButtonTxtCreateCategory)

			if err = m.setUserInfo(ctx, Send.UserID, userInfo{
				Status:  statusNonExistentCategory,
				Command: Send.Text,
			}); err != nil {
				return m.tgClient.SendMessage("Ошибочка: "+err.Error(), Send.UserID)
			}

			return m.tgClient.SendKeyboard("Такой категории у вас еще нет, выберите одну из предложенных категорий или создайте свою с помощью команды /category", Send.UserID, buttons)
		}

		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), Send.UserID)
	}

	userCur, err := m.purchasesModel.CurrencyToStr(expAndLim.Currency)
	if err != nil {
		err = errors.Wrap(err, "purchasesModel.CurrencyToStr")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), Send.UserID)
	}

	txt := ScsTxtPurchaseAdded
	if expAndLim.Limit != -1 {
		txt += fmt.Sprintf("\n\nУ вас установлен лимит: %.2f %s. За этот месяц вы потратили уже %.2f %s.",
			expAndLim.Limit, userCur, expAndLim.Expenses, userCur)
		if expAndLim.LimitExceeded {
			txt += "\nВЫ ПРЕВЫСИЛИ ЛИМИТ!"
		}
	}

	return m.tgClient.SendMessage(txt, Send.UserID)
}

func (m *Model) msgCurrency(ctx context.Context, Send Message, rawCY string) error {
	cy, err := m.purchasesModel.StrToCurrency(rawCY)
	if err != nil {
		return m.tgClient.SendMessage(ErrTxtInvalidCurrency, Send.UserID)
	}

	if err = m.purchasesModel.ChangeUserCurrency(ctx, Send.UserID, cy); err != nil {
		err = errors.Wrap(err, "purchasesModel.ChangeUserCurrency")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), Send.UserID)
	}
	return m.tgClient.SendMessage(ScsTxtCurrencyChanged, Send.UserID)
}

func (m *Model) msgLimit(ctx context.Context, Send Message, limit string) error {
	if err := m.purchasesModel.ChangeUserLimit(ctx, Send.UserID, limit); err != nil {
		if errors.Is(err, purchases.ErrLimitParsing) {
			return m.tgClient.SendMessage(ErrTxtInvalidInput, Send.UserID)
		}
		err = errors.Wrap(err, "purchasesModel.ChangeUserLimit")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), Send.UserID)
	}
	return m.tgClient.SendMessage(ScsTxtLimitChanged, Send.UserID)
}
