package purchases

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

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

type ReportItem struct {
	PurchaseCategory string
	Summa            float64
}

// Segment кусочек круговой диаграммы
type Segment struct {
	Value float64
	Label string
}

// Report создание отчета
func (m *Model) Report(period Period, userID int64) (txt string, img []byte, err error) {
	from, err := fromTime(time.Now(), period)
	if err != nil {
		return "", nil, errors.Wrap(err, "fromTime")
	}

	res, err := m.Repo.GetReport(from, userID)
	if err != nil {
		return "", nil, errors.Wrap(err, "Repo.GetReport")
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Summa > res[j].Summa
	})

	resStr := strings.Builder{}
	for _, item := range res {
		ctgr := item.PurchaseCategory
		if ctgr == "" {
			ctgr = "не указанные категории"
		}

		resStr.WriteString(ctgr)
		resStr.WriteString(": ")
		resStr.WriteString(strconv.FormatFloat(item.Summa, 'f', -1, 64))
		resStr.WriteString("\n")
	}

	resIMG, err := m.ChartDrawer.PieChart(res)
	if err != nil {
		return "", nil, errors.Wrap(err, "ChartDrawer.PieChart")
	}

	return resStr.String(), resIMG, nil
}

// fromTime позволяет получить из переданной даты новую, вычитая из переданной указанный период
// (учитывая количество дней в месяцах и високосные годы)
func fromTime(to time.Time, period Period) (time.Time, error) {
	switch period {
	case periodYear:
		from := to.AddDate(-1, 0, 0)
		days := to.Sub(from).Hours() / 24
		return to.AddDate(0, 0, int(-days)), nil
	case periodMonth:
		from := to.AddDate(0, -1, 0)
		days := to.Sub(from).Hours() / 24
		return to.AddDate(0, 0, int(-days)), nil
	case periodWeek:
		return to.AddDate(0, 0, -7), nil
	default:
		return time.Time{}, ErrUnknownPeriod
	}
}
