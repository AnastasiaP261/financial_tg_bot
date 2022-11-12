package report

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/metrics"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"strings"
	"time"
)

const keySuffix = "report"

type Request struct {
	FromDate time.Time         `json:"fromDate"`
	UserID   int64             `json:"userId"`
	Currency currency.Currency `json:"currency"`
}

type Report struct {
	Items    []ReportItem
	FromDate time.Time // дата начала выборки данных в отчете
}

type ReportItem struct {
	PurchaseCategory string
	Summa            float64
}

type Purchase struct {
	PurchaseCategory string
	Summa            float64

	// коэффициенты валют на момент совершения траты
	currency.RateToRUB
}

func (s *Service) CreateReport(ctx context.Context, rawReq string) (txt string, err error) {
	var req Request
	if err := json.Unmarshal([]byte(rawReq), &req); err != nil {
		return "", errors.Wrap(err, "unmarshalling error")
	}

	reportItems, err := s.getPurchasesReportFromDate(ctx, req.FromDate, req.UserID, req.Currency)
	if err != nil {
		return "", errors.Wrap(err, "packagingByCategory")
	}

	cy, err := currency.CurrencyToStr(req.Currency)
	if err != nil {
		return "", errors.Wrap(err, "currencyToStr")
	}

	resStr := strings.Builder{}
	resStr.WriteString("Ваша валюта: ")
	resStr.WriteString(cy)
	resStr.WriteString("\nВаш отчет:\n")
	for _, item := range reportItems {
		resStr.WriteString("\t")
		resStr.WriteString(item.PurchaseCategory)
		resStr.WriteString(": ")
		resStr.WriteString(strconv.FormatFloat(item.Summa, 'f', 2, 64))
		resStr.WriteString("\n")
	}

	return resStr.String(), nil
}

func (s *Service) getPurchasesReportFromDate(ctx context.Context, from time.Time, userID int64, cy currency.Currency) ([]ReportItem, error) {
	report, err := s.ReportsStore.GetReport(ctx, createKeyForReportsStore(userID))
	if err != nil {
		logs.Error("reports store error", zap.Error(err))
	}
	// если в хранилище статусов ничего нет или вернулась ошибка просто идем в репу
	if err == nil && len(report.Items) != 0 {
		if report.FromDate == from {
			metrics.InFlightReports.WithLabelValues(metrics.ReportSourceCache).Inc()
			return report.Items, nil
		}
	}

	purchases, err := s.Repo.GetUserPurchasesFromDate(ctx, from, userID)
	if err != nil {
		return nil, errors.Wrap(err, "Repo.GetUserPurchasesFromDate")
	}

	reportItems, err := s.packagingByCategory(purchases, cy)
	if err != nil {
		return nil, errors.Wrap(err, "packagingByCategory")
	}

	err = s.ReportsStore.SetReport(ctx, createKeyForReportsStore(userID), Report{Items: reportItems, FromDate: from}) // nolint: errcheck
	if err != nil {
		return nil, errors.Wrap(err, "ReportsStore.SetReport")
	}
	metrics.InFlightReports.WithLabelValues(metrics.ReportSourceBD).Inc()

	logs.Info("REPORT", zap.Any("items", reportItems))

	return reportItems, nil
}

// packagingByCategory получает на вход список трат и формирует из него отчет, переводя все траты в
// выбранную валюту и складывая их по категориям
func (s *Service) packagingByCategory(purchases []Purchase, currentCurrency currency.Currency) ([]ReportItem, error) {
	tempCategoryOnSum := make(map[string]float64, len(purchases))
	for _, p := range purchases {
		resSum, err := currency.RubToCurrentCurrency(currentCurrency, p.Summa, p.RateToRUB)
		if err != nil {
			return nil, errors.Wrap(err, "rubToCurrentCurrency")
		}

		tempCategoryOnSum[p.PurchaseCategory] += resSum
	}

	res := make([]ReportItem, 0, len(tempCategoryOnSum))
	for k, v := range tempCategoryOnSum {
		res = append(res, ReportItem{
			PurchaseCategory: k,
			Summa:            v,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Summa > res[j].Summa
	})

	return res, nil
}

func createKeyForReportsStore(userID int64) string {
	return strconv.FormatInt(userID, 10) + keySuffix
}
