package purchases

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/wrappers/metrics"
	"go.uber.org/zap"
)

const keySuffix = "report"

type Period byte

var (
	periodYear  Period = 1
	periodMonth Period = 2
	periodWeek  Period = 3
)

func (m *Model) ToPeriod(str string) (Period, error) {
	switch str {
	case "year":
		return periodYear, nil
	case "month":
		return periodMonth, nil
	case "week":
		return periodWeek, nil
	default:
		return 0, ErrUnknownPeriod
	}
}

type Purchase struct {
	PurchaseCategory string
	Summa            float64

	// коэффициенты валют на момент совершения траты
	currency.RateToRUB
}

type ReportItem struct {
	PurchaseCategory string
	Summa            float64
}

type Report struct {
	Items    []ReportItem
	FromDate time.Time // дата начала выборки данных в отчете
}

// Segment кусочек круговой диаграммы
type Segment struct {
	Value float64
	Label string
}

type ReportRequest struct {
	FromDate time.Time         `json:"fromDate"`
	UserID   int64             `json:"userId"`
	Currency currency.Currency `json:"currency"`
}

// Report создание отчета
func (m *Model) Report(ctx context.Context, period Period, userID int64) (txt string, img []byte, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "report")
	defer span.Finish()

	from, err := fromTime(time.Now(), period)
	if err != nil {
		return "", nil, errors.Wrap(err, "fromTime")
	}

	info, err := m.Repo.GetUserInfo(ctx, userID)
	if err != nil {
		return "", nil, errors.Wrap(err, "Repo.GetUserInfo")
	}

	{
		jsonReq, err := json.Marshal(ReportRequest{
			FromDate: from,
			UserID:   userID,
			Currency: info.Currency,
		})
		if err != nil {
			return "", nil, errors.Wrap(err, "marshalling error")
		}
		_ = m.BrokerMsgCreator.SendNewMsg(
			"get_report",
			string(jsonReq),
		)
	}

	reportItems, err := m.getPurchasesReportFromDate(ctx, from, userID, info)
	if err != nil {
		return "", nil, errors.Wrap(err, "packagingByCategory")
	}

	cy, err := currency.CurrencyToStr(info.Currency)
	if err != nil {
		return "", nil, errors.Wrap(err, "CurrencyToStr")
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

	resIMG, err := m.ChartDrawer.PieChart(reportItems)
	if err != nil {
		return "", nil, errors.Wrap(err, "ChartDrawer.PieChart")
	}

	return resStr.String(), resIMG, nil
}

func (m *Model) getPurchasesReportFromDate(ctx context.Context, from time.Time, userID int64, info User) ([]ReportItem, error) {
	report, err := m.ReportsStore.GetReport(ctx, createKeyForReportsStore(userID))
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

	purchases, err := m.Repo.GetUserPurchasesFromDate(ctx, from, userID)
	if err != nil {
		return nil, errors.Wrap(err, "Repo.GetUserPurchasesFromDate")
	}

	reportItems, err := m.packagingByCategory(purchases, info.Currency)
	if err != nil {
		return nil, errors.Wrap(err, "packagingByCategory")
	}

	m.ReportsStore.SetReport(ctx, createKeyForReportsStore(userID), Report{Items: reportItems, FromDate: from}) // nolint: errcheck
	metrics.InFlightReports.WithLabelValues(metrics.ReportSourceBD).Inc()

	return reportItems, nil
}

// packagingByCategory получает на вход список трат и формирует из него отчет, переводя все траты в
// выбранную валюту и складывая их по категориям
func (m *Model) packagingByCategory(purchases []Purchase, currentCurrency currency.Currency) ([]ReportItem, error) {
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

// fromTime позволяет получить из переданной даты новую, вычитая из переданной указанный период
// (учитывая количество дней в месяцах и високосные годы)
func fromTime(to time.Time, period Period) (time.Time, error) {
	var resDate time.Time

	switch period {
	case periodYear:
		from := to.AddDate(-1, 0, 0)
		days := to.Sub(from).Hours() / 24
		resDate = to.AddDate(0, 0, int(-days))

	case periodMonth:
		from := to.AddDate(0, -1, 0)
		days := to.Sub(from).Hours() / 24
		resDate = to.AddDate(0, 0, int(-days))

	case periodWeek:
		resDate = to.AddDate(0, 0, -7)

	default:
		return time.Time{}, ErrUnknownPeriod
	}

	// этот шаг нужен, чтобы обнулить все составные части кроме даты
	y, m, d := resDate.Date()
	result, err := time.Parse("02.01.2006", fmt.Sprintf("%02d.%02d.%d", d, m, y))
	if err != nil {
		return time.Time{}, errors.Wrap(err, "time.Parse")
	}

	return result, nil
}

func createKeyForReportsStore(userID int64) string {
	return strconv.FormatInt(userID, 10) + keySuffix
}
