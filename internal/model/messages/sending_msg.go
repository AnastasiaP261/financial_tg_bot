package messages

import "context"

type ReportData struct {
	UserID        int64
	ReportMessage string
	ReportIMG     []byte
}

func (m *Model) SendReport(ctx context.Context, userID int64, text string, img []byte) error {
	_ = m.tgClient.SendMessage(ScsTxtCategoryCreated, userID)
	_ = m.tgClient.SendMessage(text, userID)
	return m.tgClient.SendImage(img, userID)
}
