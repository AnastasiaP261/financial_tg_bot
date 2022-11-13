package reports

import (
	"context"

	"github.com/pkg/errors"
	pkg "gitlab.ozon.dev/apetrichuk/financial-tg-bot/pkg/api/financial-tg-bot"
)

func (s *server) SendReport(ctx context.Context, req *pkg.SendReportRequest) (*pkg.SendReportResponse, error) {
	if err := s.sender.SendReport(ctx, req.GetUserId(), req.GetReportMessage(), req.GetReportImage()); err != nil {
		return nil, errors.Wrap(err, "sender.SendReport")
	}

	return &pkg.SendReportResponse{Response: &pkg.DefaultResponse{
		Success:      true,
		ErrorMessage: "",
	}}, nil
}
