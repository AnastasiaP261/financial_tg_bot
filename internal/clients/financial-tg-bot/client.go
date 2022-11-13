package financial_tg_bot

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/report"
	pb "gitlab.ozon.dev/apetrichuk/financial-tg-bot/pkg/api/financial-tg-bot"
)

type service struct {
	pb.MessagesServiceClient
}

func New() *service {
	return &service{}
}

func (s *service) SendReport(ctx context.Context, req report.SendReportRequest) (report.SendReportResponse, error) {
	res, err := s.MessagesServiceClient.SendReport(ctx, &pb.SendReportRequest{
		UserId:        req.UserId,
		ReportMessage: req.ReportMessage,
	})
	if err != nil {
		return report.SendReportResponse{}, errors.Wrap(err, "MessagesServiceClient.SendReport")
	}

	return report.SendReportResponse{
		Response: report.DefaultResponse{
			Success:      res.GetResponse().GetSuccess(),
			ErrorMessage: res.GetResponse().GetErrorMessage(),
		},
	}, nil
}
