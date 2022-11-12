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

type Service struct {
	Repo         Repo
	ReportsStore ReportsStore
}
