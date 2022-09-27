package messages

import (
	"regexp"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type PurchasesModel interface {
	AddPurchase(userID int64, rawSum, category, rawDate string) error
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
	Text   string
	UserID int64
}

var (
	addConditionOnlySum               = regexp.MustCompile(`/add (\d+.?\d*)`)
	addConditionSumAndCategory        = regexp.MustCompile(`/add (\d+.?\d*) ([ \wФА-Яа-я]+)`)
	addConditionSumAndCategoryAndDate = regexp.MustCompile(`/add (\d+\.?\d*) ([ \wФА-Яа-я]+) (\d{2}\.\d{2}\.\d{4})`)
)

func (s *Model) IncomingMessage(msg Message) error {
	switch {
	case msg.Text == "/start":
		return s.tgClient.SendMessage("hello", msg.UserID)

	case addConditionSumAndCategoryAndDate.MatchString(msg.Text):
		res := addConditionSumAndCategoryAndDate.FindStringSubmatch(msg.Text)
		if len(res) < 4 {
			return s.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID)
		}

		err := s.purchasesModel.AddPurchase(msg.UserID, res[1], res[2], res[3])
		if err != nil {
			return s.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID)
		}
		return s.tgClient.SendMessage(ScsTxtPurchaseAdded, msg.UserID)

	case addConditionSumAndCategory.MatchString(msg.Text):
		res := addConditionSumAndCategory.FindStringSubmatch(msg.Text)
		if len(res) < 3 {
			return s.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID)
		}

		err := s.purchasesModel.AddPurchase(msg.UserID, res[1], res[2], "")
		if err != nil {
			return s.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID)
		}
		return s.tgClient.SendMessage(ScsTxtPurchaseAdded, msg.UserID)

	case addConditionOnlySum.MatchString(msg.Text):
		res := addConditionOnlySum.FindStringSubmatch(msg.Text)
		if len(res) < 2 {
			return s.tgClient.SendMessage(ErrTxtInvalidInput, msg.UserID)
		}

		err := s.purchasesModel.AddPurchase(msg.UserID, res[1], "", "")
		if err != nil {
			return s.tgClient.SendMessage("Ошибочка: "+err.Error(), msg.UserID)
		}
		return s.tgClient.SendMessage(ScsTxtPurchaseAdded, msg.UserID)

	default:
		return s.tgClient.SendMessage(ErrTxtUnknownCommand, msg.UserID)
	}
}
