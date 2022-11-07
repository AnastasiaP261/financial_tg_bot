package messages

import (
	"context"
	"regexp"
)

type Message struct {
	Text     string
	UserID   int64
	UserName string
}

var (
	// addPurchaseOnlySum сообщение о добавлении траты без категории и даты (указывается текущая дата)
	addPurchaseOnlySum = regexp.MustCompile(`/add (\d+.?\d*)`)
	// addPurchaseSumAndCategory сообщение о добавлении траты с категорией, но без даты (указывается текущая дата)
	addPurchaseSumAndCategory = regexp.MustCompile(`/add (\d+.?\d*) ([ \wФА-Яа-я]+)`)
	// addPurchaseSumAndCategoryAndDate сообщение о добавлении траты с категорией и датой
	addPurchaseSumAndCategoryAndDate = regexp.MustCompile(`/add (\d+\.?\d*) ([ \wФА-Яа-я]+) (\d{2}\.\d{2}\.\d{4})`)

	// addCategory добавление новой категории
	addCategory = regexp.MustCompile(`/category ([ \wФА-Яа-я\-]+)`)

	// report создание отчета за выбранный период
	report = regexp.MustCompile(`/report (month|week|year)`)

	// currency команда для смены основной валюты пользователя
	currency = regexp.MustCompile(`/currency ([A-Za-z]{3})`)
	// limit команда для задания месячного лимита трат пользователю
	limit = regexp.MustCompile(`/limit (\d+.?\d*)`)
)

func (m *Model) IncomingMessage(ctx context.Context, msg Message) error {
	switch {
	case msg.Text == "/start":
		return m.SendMessage("hello", msg.UserID)

	case report.MatchString(msg.Text):
		return m.msgReport(ctx, msg)

	case addCategory.MatchString(msg.Text):
		return m.msgAddCategory(ctx, msg)

	case addPurchaseSumAndCategoryAndDate.MatchString(msg.Text):
		res := addPurchaseSumAndCategoryAndDate.FindStringSubmatch(msg.Text)
		if len(res) < 4 {
			return m.SendMessage(ErrTxtInvalidInput, msg.UserID)
		}

		return m.msgAddPurchase(ctx, msg, res[1], res[2], res[3])

	case addPurchaseSumAndCategory.MatchString(msg.Text):
		res := addPurchaseSumAndCategory.FindStringSubmatch(msg.Text)
		if len(res) < 3 {
			return m.SendMessage(ErrTxtInvalidInput, msg.UserID)
		}

		return m.msgAddPurchase(ctx, msg, res[1], res[2], "")

	case addPurchaseOnlySum.MatchString(msg.Text):
		res := addPurchaseOnlySum.FindStringSubmatch(msg.Text)
		if len(res) < 2 {
			return m.SendMessage(ErrTxtInvalidInput, msg.UserID)
		}

		return m.msgAddPurchase(ctx, msg, res[1], "", "")

	case currency.MatchString(msg.Text):
		res := currency.FindStringSubmatch(msg.Text)
		if len(res) < 2 {
			return m.SendMessage(ErrTxtInvalidInput, msg.UserID)
		}

		return m.msgCurrency(ctx, msg, res[1])

	case limit.MatchString(msg.Text):
		res := limit.FindStringSubmatch(msg.Text)
		if len(res) < 2 {
			return m.SendMessage(ErrTxtInvalidInput, msg.UserID)
		}

		return m.msgLimit(ctx, msg, res[1])

	default:
		return m.SendMessage(ErrTxtUnknownCommand, msg.UserID)
	}
}
