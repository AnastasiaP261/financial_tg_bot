package messages

import (
	"regexp"

	"github.com/pkg/errors"
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
	// addConditionOnlySum сообщение о добавлении траты без категории и даты (указывается текущая дата)
	addConditionOnlySum = regexp.MustCompile(`/add (\d+.?\d*)`)
	// addConditionSumAndCategory сообщение о добавлении траты c категорией но без даты (указывается текущая дата)
	addConditionSumAndCategory = regexp.MustCompile(`/add (\d+.?\d*) ([ \wФА-Яа-я]+)`)
	// addConditionSumAndCategoryAndDate сообщение о добавлении траты c категорией и датой
	addConditionSumAndCategoryAndDate = regexp.MustCompile(`/add (\d+\.?\d*) ([ \wФА-Яа-я]+) (\d{2}\.\d{2}\.\d{4})`)

	// addCategory добавление новой категории
	addCategory = regexp.MustCompile(`/category ([ \wФА-Яа-я\-]+)`)

	// report создание отчета за выбранный период
	report = regexp.MustCompile(`/report (month|week|year)`)
)

func (s *Model) IncomingMessage(msg Message) error {
	switch {
	case msg.Text == "/start":
		return s.tgClient.SendMessage("hello", msg.UserID, msg.UserName)

	case report.MatchString(msg.Text):
		res := report.FindStringSubmatch(msg.Text)
		if len(res) < 2 {
			return s.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}

		period, err := s.purchasesModel.ToPeriod(res[1])
		if err != nil {
			return s.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}

		txt, img, err := s.purchasesModel.Report(period, msg.UserID)
		if err != nil {
			return s.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
		}

		if err = s.tgClient.SendMessage("Ваш отчет:\n\n"+txt, msg.UserID, msg.UserName); err != nil {
			return err
		}

		return s.tgClient.SendImage(img, msg.ChatID, msg.UserName)

	case addCategory.MatchString(msg.Text):
		res := addCategory.FindStringSubmatch(msg.Text)
		if len(res) < 2 {
			return s.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}

		err := s.purchasesModel.AddCategory(msg.UserID, res[1])
		if err != nil {
			return s.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
		}
		return s.tgClient.SendMessage(ScsTxtCategoryAdded, msg.UserID, msg.UserName)

	case addConditionSumAndCategoryAndDate.MatchString(msg.Text):
		res := addConditionSumAndCategoryAndDate.FindStringSubmatch(msg.Text)
		if len(res) < 4 {
			return s.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}

		err := s.purchasesModel.AddPurchase(msg.UserID, res[1], res[2], res[3])
		if err != nil {
			if errors.Is(err, purchases.ErrCategoryNotExist) {
				return s.tgClient.SendMessage(ErrTxtCategoryDoesntExist, msg.UserID, msg.UserName)
			}
			return s.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
		}
		return s.tgClient.SendMessage(ScsTxtPurchaseAdded, msg.UserID, msg.UserName)

	case addConditionSumAndCategory.MatchString(msg.Text):
		res := addConditionSumAndCategory.FindStringSubmatch(msg.Text)
		if len(res) < 3 {
			return s.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}

		err := s.purchasesModel.AddPurchase(msg.UserID, res[1], res[2], "")
		if err != nil {
			if errors.Is(err, purchases.ErrCategoryNotExist) {
				return s.tgClient.SendMessage(ErrTxtCategoryDoesntExist, msg.UserID, msg.UserName)
			}
			return s.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
		}
		return s.tgClient.SendMessage(ScsTxtPurchaseAdded, msg.UserID, msg.UserName)

	case addConditionOnlySum.MatchString(msg.Text):
		res := addConditionOnlySum.FindStringSubmatch(msg.Text)
		if len(res) < 2 {
			return s.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID, msg.UserName)
		}

		err := s.purchasesModel.AddPurchase(msg.UserID, res[1], "", "")
		if err != nil {
			if errors.Is(err, purchases.ErrCategoryNotExist) {
				return s.tgClient.SendMessage(ErrTxtCategoryDoesntExist, msg.UserID, msg.UserName)
			}
			return s.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID, msg.UserName)
		}
		return s.tgClient.SendMessage(ScsTxtPurchaseAdded, msg.UserID, msg.UserName)

	default:
		return s.tgClient.SendMessage(ErrTxtUnknownCommand, msg.UserID, msg.UserName)
	}
}
