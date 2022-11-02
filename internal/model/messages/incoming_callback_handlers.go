package messages

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/normalize"
)

func (m *Model) msgNonExistentCategory(ctx context.Context, msg Callback, info userInfo) error {
	// если пользователь выбрал создание новой категории
	if msg.Data == ButtonTxtCreateCategory {
		if err := m.setUserInfo(ctx, msg.UserID, userInfo{}); err != nil {
			err = errors.Wrap(err, "setUserInfo")
			return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
		}
		return m.tgClient.SendMessage(ScsTxtCategoryAddSelected, msg.UserID, msg.UserName)
	}
	catName := msg.Data

	// если пользователь выбрал одну из предложенных категорий
	userHasThisCat := false
	{
		userCats, err := m.purchasesModel.GetUserCategories(ctx, msg.UserID)
		if err != nil {
			return errors.Wrap(err, "purchasesModel.GetUserCategories")
		}
		for _, c := range userCats {
			if catName == c {
				userHasThisCat = true
				break
			}
		}
	}

	// если он выбрал категорию которая у него еще не добавлена
	if !userHasThisCat {
		if err := m.purchasesModel.AddCategoryToUser(ctx, msg.UserID, normalize.Category(catName)); err != nil {
			err = errors.Wrap(err, "purchasesModel.AddCategoryToUser")
			return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
		}
		_ = m.tgClient.SendMessage(ScsTxtCategoryAddedToUser, msg.UserID, msg.UserName)
	}

	if err := m.setUserInfo(ctx, msg.UserID, userInfo{}); err != nil {
		err = errors.Wrap(err, "setUserInfo")
		return m.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
	}

	return m.IncomingMessage(ctx, Message{
		Text:     info.Command,
		UserID:   msg.UserID,
		UserName: msg.UserName,
	})
}
