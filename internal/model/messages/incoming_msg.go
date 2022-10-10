package messages

import (
	"regexp"

	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

type MessageSender interface {
	SendMessage(text string, userID int64, userName string) error
	SendImage(img []byte, chatID int64, userName string) error
}

type PurchasesModel interface {
	AddPurchase(userID int64, rawSum, category, rawDate string) error
	AddCategory(userID int64, category string) error
	Report(period purchases.Period, userID int64) (txt string, img []byte, err error)
	ToPeriod(str string) (purchases.Period, error)
}

type Model struct {
	tgClient       MessageSender
	purchasesModel PurchasesModel
}

func New(tgClient MessageSender, purchasesModel PurchasesModel) *Model {
	return &Model{
		tgClient:       tgClient,
		purchasesModel: purchasesModel,
	}
}

type Message struct {
	Text     string
	UserID   int64
	ChatID   int64
	UserName string
}

var (
	// addPurchaseOnlySum сообщение о добавлении траты без категории и даты (указывается текущая дата)
	addPurchaseOnlySum = regexp.MustCompile(`/add (\d+.?\d*)`)
	// addPurchaseSumAndCategory сообщение о добавлении траты c категорией но без даты (указывается текущая дата)
	addPurchaseSumAndCategory = regexp.MustCompile(`/add (\d+.?\d*) ([ \wФА-Яа-я]+)`)
	// addPurchaseSumAndCategoryAndDate сообщение о добавлении траты c категорией и датой
	addPurchaseSumAndCategoryAndDate = regexp.MustCompile(`/add (\d+\.?\d*) ([ \wФА-Яа-я]+) (\d{2}\.\d{2}\.\d{4})`)

	// addCategory добавление новой категории
	addCategory = regexp.MustCompile(`/category ([ \wФА-Яа-я\-]+)`)

	// report создание отчета за выбранный период
	report = regexp.MustCompile(`/report (month|week|year)`)
)

func (m *Model) IncomingMessage(msg Message) error {
	switch {
	case msg.Text == "/start":
		return m.tgClient.SendMessage("hello", msg.UserID, msg.UserName)

	case report.MatchString(msg.Text):
		return m.msgReport(msg)

	case addCategory.MatchString(msg.Text):
		return m.msgAddCategory(msg)

	case addPurchaseSumAndCategoryAndDate.MatchString(msg.Text):
		res := addPurchaseSumAndCategoryAndDate.FindStringSubmatch(msg.Text)
		if len(res) < 4 {
			return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}

		return m.msgAddPurchase(msg, res[1], res[2], res[3])

	case addPurchaseSumAndCategory.MatchString(msg.Text):
		res := addPurchaseSumAndCategory.FindStringSubmatch(msg.Text)
		if len(res) < 3 {
			return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}

		return m.msgAddPurchase(msg, res[1], res[2], "")

	case addPurchaseOnlySum.MatchString(msg.Text):
		res := addPurchaseOnlySum.FindStringSubmatch(msg.Text)
		if len(res) < 2 {
			return m.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}

		return m.msgAddPurchase(msg, res[1], "", "")

	default:
		return m.tgClient.SendMessage(ErrTxtUnknownCommand, msg.UserID, msg.UserName)
	}
}
