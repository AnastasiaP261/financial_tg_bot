package financial_reports

import (
	"context"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/report"
)

type ReportHandler interface {
	CreateReport(ctx context.Context, rawReq string) (report.CreateReportResponse, error)
}
