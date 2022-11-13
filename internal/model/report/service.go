package report

import (
	"context"
	"time"
)

// Repo репозиторий
type Repo interface {
	GetUserPurchasesFromDate(ctx context.Context, fromDate time.Time, userID int64) ([]Purchase, error)
}

type ReportsStore interface {
	SetReport(ctx context.Context, key string, value Report) error
	GetReport(ctx context.Context, key string) (Report, error)
	Delete(ctx context.Context, key string) error
}

type Sender interface {
	SendReport(ctx context.Context, req SendReportRequest) (SendReportResponse, error)
}

type service struct {
	repo         Repo
	reportsStore ReportsStore
	sender       Sender
}

func New(repo Repo, store ReportsStore, sender Sender) *service {
	return &service{
		repo:         repo,
		reportsStore: store,
		sender:       sender,
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
