package messages

import "context"

type Callback struct {
	UserID   int64
	UserName string
	Data     string
}

func (m *Model) IncomingCallback(ctx context.Context, msg Callback) error {
	info, err := m.getUserInfo(ctx, msg.UserID)
	if err != nil {
		return m.SendMessage("Ошибочка: "+err.Error(), msg.UserID)
	}

	switch info.Status {
	case statusNonExistentCategory:
		return m.msgNonExistentCategory(ctx, msg, info)

	default:
		if err = m.setUserInfo(ctx, msg.UserID, userInfo{}); err != nil {
			return m.SendMessage("Ошибочка: "+err.Error(), msg.UserID)
		}
		return m.SendMessage(ErrTxtInvalidStatus, msg.UserID)
	}
}
