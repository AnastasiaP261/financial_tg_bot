package reports

import (
	"context"
	"fmt"
	pkg "gitlab.ozon.dev/apetrichuk/financial-tg-bot/pkg/api/financial-tg-bot"
)

func (s *server) SendReport(ctx context.Context, req *pkg.SendReportRequest) (*pkg.SendReportResponse, error) {
	fmt.Printf("### TG-BOT принято сообщение: %-v", req)

	return &pkg.SendReportResponse{Response: &pkg.DefaultResponse{
		Success:      true,
		ErrorMessage: "",
	}}, nil
}
