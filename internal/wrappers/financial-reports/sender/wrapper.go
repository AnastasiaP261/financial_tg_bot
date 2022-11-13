package sender

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/report"
	financial_reports "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/financial-reports"
)

type Sender interface {
	SendReport(ctx context.Context, req report.SendReportRequest) (report.SendReportResponse, error)
}

type Wrapper struct {
	handler financial_reports.ReportHandler
	client  Sender
}

func NewWrapper(handler financial_reports.ReportHandler, client Sender) *Wrapper {
	return &Wrapper{
		handler: handler,
		client:  client,
	}
}

func (w *Wrapper) CreateReport(ctx context.Context, rawReq string) (report.CreateReportResponse, error) {
	res, err := w.handler.CreateReport(ctx, rawReq)
	if err != nil {
		return report.CreateReportResponse{}, err
	}

	resp, err := w.client.SendReport(ctx, report.SendReportRequest{
		UserID:        res.UserID,
		ReportMessage: res.Text,
		ReportIMG:     res.IMG,
	})
	if err != nil {
		return report.CreateReportResponse{}, errors.Wrap(err, "sender.SendReport")
	}
	if !resp.Response.Success {
		return report.CreateReportResponse{}, errors.New("report send failed")
	}

	return res, err
}
