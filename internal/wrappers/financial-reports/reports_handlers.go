package financial_reports

import "context"

type ReportHandler interface {
	CreateReport(ctx context.Context, rawReq string) (string, int64, error)
}
