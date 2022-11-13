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

// ChartDrawer рисовальщик
type ChartDrawer interface {
	// PieChart нарисовать круговую диаграмму трат
	PieChart(data []ReportItem) ([]byte, error)
}

type service struct {
	repo         Repo
	reportsStore ReportsStore
	Drawer       ChartDrawer
}

func New(repo Repo, store ReportsStore, drawer ChartDrawer) *service {
	return &service{
		repo:         repo,
		reportsStore: store,
		Drawer:       drawer,
	}
}

type SendReportRequest struct {
	UserId        int64
	ReportMessage string
	ReportIMG     []byte
}

type SendReportResponse struct {
	Response DefaultResponse
}

type DefaultResponse struct {
	Success      bool
	ErrorMessage string
}
