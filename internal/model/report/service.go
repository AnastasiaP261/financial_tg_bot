package report

import (
	"context"
	"time"

	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

// Repo репозиторий
type Repo interface {
	GetUserPurchasesFromDate(ctx context.Context, fromDate time.Time, userID int64) ([]purchases.Purchase, error)
}

type ReportsStore interface {
	SetReport(ctx context.Context, key string, value Report) error
	GetReport(ctx context.Context, key string) (Report, error)
	Delete(ctx context.Context, key string) error
}

type service struct {
	repo         Repo
	reportsStore ReportsStore
}

func New(repo Repo, store ReportsStore) *service {
	return &service{
		repo:         repo,
		reportsStore: store,
	}
}

type SendReportRequest struct {
	UserId        int64
	ReportMessage string
}

type SendReportResponse struct {
	Response DefaultResponse
}

type DefaultResponse struct {
	Success      bool
	ErrorMessage string
}
