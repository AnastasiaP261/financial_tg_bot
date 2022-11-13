package financial_tg_bot

import (
	"context"
	"fmt"

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
	GRPCPortMessages() int
}

func New(confService configGetter) (*service, error) {
	address := fmt.Sprintf("%s:%d", confService.GRPCHostMessages(), confService.GRPCPortMessages())
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &service{cl: pb.NewMessagesServiceClient(conn)}, nil
}

func (s *service) SendReport(ctx context.Context, req report.SendReportRequest) (report.SendReportResponse, error) {
	res, err := s.cl.SendReport(ctx, &pb.SendReportRequest{
		UserId:        req.UserId,
		ReportMessage: req.ReportMessage,
		ReportImage:   req.ReportIMG,
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
