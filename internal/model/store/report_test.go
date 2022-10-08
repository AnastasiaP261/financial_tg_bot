package store

import (
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	"testing"
	"time"
)

func Test_GetReport(t *testing.T) {
	t.Run("отчет в рублях", func(t *testing.T) {
		s := New()

		t1, _ := time.Parse("02.01.2006", "01.10.2022")
		t2, _ := time.Parse("02.01.2006", "02.10.2022")
		t3, _ := time.Parse("02.01.2006", "03.10.2022")
		t4, _ := time.Parse("02.01.2006", "04.10.2022")

		s.Users = []user{
			{UserID: 123, Currency: RUB},
		}
		s.Purchases = []purchase{
			{UserID: 123, Sum: 100, Category: "some category", Date: t1, USDRatio: 2, CNYRatio: 2, EURRatio: 2},
			{UserID: 123, Sum: 100, Category: "some category", Date: t2, USDRatio: 2, CNYRatio: 2, EURRatio: 2},
			{UserID: 123, Sum: 100, Category: "some category", Date: t3, USDRatio: 2, CNYRatio: 2, EURRatio: 2},
			{UserID: 123, Sum: 100, Category: "some category", Date: t4, USDRatio: 2, CNYRatio: 2, EURRatio: 2},
		}

		res, err := s.GetReport(t2, 123)

		assert.NoError(t, err)
		assert.Equal(t, []purchases.ReportItem{
			{"some category", 100 + 100 + 100},
		}, res)
	})

	t.Run("отчет в валюте", func(t *testing.T) {
		s := New()

		t1, _ := time.Parse("02.01.2006", "01.10.2022")
		t2, _ := time.Parse("02.01.2006", "02.10.2022")
		t3, _ := time.Parse("02.01.2006", "03.10.2022")
		t4, _ := time.Parse("02.01.2006", "04.10.2022")

		s.Users = []user{
			{UserID: 123, Currency: USD},
		}
		s.Purchases = []purchase{
			{UserID: 123, Sum: 100, Category: "some category", Date: t1, USDRatio: 1.5, CNYRatio: 2, EURRatio: 2},
			{UserID: 123, Sum: 100, Category: "some category", Date: t2, USDRatio: 1.5, CNYRatio: 2, EURRatio: 2},
			{UserID: 123, Sum: 100, Category: "some category", Date: t3, USDRatio: 1.5, CNYRatio: 2, EURRatio: 2},
			{UserID: 123, Sum: 100, Category: "some category", Date: t4, USDRatio: 1.7, CNYRatio: 2, EURRatio: 2},
		}

		res, err := s.GetReport(t2, 123)

		assert.NoError(t, err)
		assert.Equal(t, []purchases.ReportItem{
			{"some category", 100*1.5 + 100*1.5 + 100*1.7},
		}, res)
	})
}
