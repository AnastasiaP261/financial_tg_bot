package model

type MessageSender interface {
	SendMessage(userID uint64, text string) error
}

type Model struct {
	TGClient MessageSender
}

func New(tgCl MessageSender) *Model {
	return &Model{TGClient: tgCl}
}

type Message struct {
	UserID uint64
	Text   string
}

func (m *Model) IncomingMsg(msg Message) {

}
