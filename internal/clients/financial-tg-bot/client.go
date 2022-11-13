package financial_tg_bot

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/report"
	pb "gitlab.ozon.dev/apetrichuk/financial-tg-bot/pkg/api/financial-tg-bot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type service struct {
	cl pb.MessagesServiceClient
}

type configGetter interface {
	GRPCHostMessages() string
}

func New(confService configGetter) (*service, error) {
	conn, err := grpc.Dial(confService.GRPCHostMessages(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &service{cl: pb.NewMessagesServiceClient(conn)}, nil
}

func (s *service) SendReport(ctx context.Context, req report.SendReportRequest) (report.SendReportResponse, error) {
	res, err := s.cl.SendReport(ctx, &pb.SendReportRequest{
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
