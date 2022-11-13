package purchases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/currency"
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

// CreateReportRequest создание запроса на отчет
func (m *Model) CreateReportRequest(ctx context.Context, period Period, userID int64) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "report")
	defer span.Finish()

	from, err := fromTime(time.Now(), period)
	if err != nil {
		return errors.Wrap(err, "fromTime")
	}

	info, err := m.Repo.GetUserInfo(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "repo.GetUserInfo")
	}

	jsonReq, err := json.Marshal(ReportRequest{
		FromDate: from,
		UserID:   userID,
		Currency: info.Currency,
	})
	if err != nil {
		return errors.Wrap(err, "marshalling error")
	}

	if err = m.BrokerMsgCreator.SendNewMsg(
		"get_report",
		string(jsonReq),
	); err != nil {
		return ErrCreateReportRequest
	}

	return nil
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
