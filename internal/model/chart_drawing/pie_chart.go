package chart_drawing

import (
	"bytes"

	"github.com/pkg/errors"
	chart "github.com/wcharczuk/go-chart/v2"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

type Model struct{}

func New() *Model {
	return &Model{}
}

// PieChart генерирует круговую диаграмму, на которой указаны траты по категориям
func (m *Model) PieChart(data []model.ReportItem) ([]byte, error) {
	values := make([]chart.Value, len(data))
	for i := range values {
		values[i] = chart.Value{
			Value: data[i].Summa,
			Label: data[i].PurchaseCategory,
		}
	}

	pie := chart.PieChart{
		Width:  1000,
		Height: 1000,
		Values: values,
	}

	img := bytes.NewBuffer([]byte{})

	err := pie.Render(chart.PNG, img)
	if err != nil {
		return nil, errors.Wrap(err, "pie.Render")
	}

	return img.Bytes(), nil
}
