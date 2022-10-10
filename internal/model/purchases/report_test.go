package purchases

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_packagingByCategory(t *testing.T) {
	m := New(nil, nil, nil)

	res, err := m.packagingByCategory([]Purchase{
		{PurchaseCategory: "cat1", Summa: 100, RateToRUB: RateToRUB{
			USD: 2,
			EUR: 2,
			CNY: 2,
		}},
		{PurchaseCategory: "cat2", Summa: 150, RateToRUB: RateToRUB{
			USD: 2,
			EUR: 2,
			CNY: 2,
		}},
		{PurchaseCategory: "", Summa: 100, RateToRUB: RateToRUB{
			USD: 2,
			EUR: 2,
			CNY: 2,
		}},
		{PurchaseCategory: "cat1", Summa: 50, RateToRUB: RateToRUB{
			USD: 2,
			EUR: 2,
			CNY: 2,
		}},
		{PurchaseCategory: "cat3", Summa: 350, RateToRUB: RateToRUB{
			USD: 2,
			EUR: 2,
			CNY: 2,
		}},
		{PurchaseCategory: "cat1", Summa: 120, RateToRUB: RateToRUB{
			USD: 2,
			EUR: 2,
			CNY: 2,
		}},
		{PurchaseCategory: "", Summa: 200, RateToRUB: RateToRUB{
			USD: 2,
			EUR: 2,
			CNY: 2,
		}},
	}, RUB)

	assert.NoError(t, err)
	assert.Equal(t, []ReportItem{
		{PurchaseCategory: "cat3", Summa: 350},
		{PurchaseCategory: "не указанные категории", Summa: 100 + 200},
		{PurchaseCategory: "cat1", Summa: 100 + 50 + 120},
		{PurchaseCategory: "cat2", Summa: 150},
	}, res)
}

func Test_fromTime(t *testing.T) {
	t.Run("неделя", func(t *testing.T) {
		in, _ := time.Parse("02.01.2006", "01.10.2022")
		out, _ := time.Parse("02.01.2006", "24.09.2022")

		res, err := fromTime(in, periodWeek)

		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})

	t.Run("месяц; у текущего 31 день у предыдущего 30", func(t *testing.T) {
		in, _ := time.Parse("02.01.2006", "31.10.2022")
		out, _ := time.Parse("02.01.2006", "01.10.2022")

		res, err := fromTime(in, periodMonth)

		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})

	t.Run("месяц; у текущего 30 дней у предыдущего 31", func(t *testing.T) {
		in, _ := time.Parse("02.01.2006", "30.09.2022")
		out, _ := time.Parse("02.01.2006", "30.08.2022")

		res, err := fromTime(in, periodMonth)

		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})

	t.Run("год; вернется та же дата, если предыдущий год високосный (дата после добавочного дня)", func(t *testing.T) {
		in, _ := time.Parse("02.01.2006", "30.09.2021")
		out, _ := time.Parse("02.01.2006", "30.09.2020")

		res, err := fromTime(in, periodYear)

		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})

	t.Run("год; вернется та же дата, если текущий год високосный (дата после добавочного дня)", func(t *testing.T) {
		in, _ := time.Parse("02.01.2006", "30.09.2020")
		out, _ := time.Parse("02.01.2006", "30.09.2019")

		res, err := fromTime(in, periodYear)

		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})

	t.Run("год; вернется та же дата, если предыдущий год високосный (дата перед добавочным днем)", func(t *testing.T) {
		in, _ := time.Parse("02.01.2006", "31.01.2021")
		out, _ := time.Parse("02.01.2006", "31.01.2020")

		res, err := fromTime(in, periodYear)

		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})

	t.Run("год; вернется та же дата, если текущий год високосный (дата перед добавочным днем)", func(t *testing.T) {
		in, _ := time.Parse("02.01.2006", "31.01.2020")
		out, _ := time.Parse("02.01.2006", "31.01.2019")

		res, err := fromTime(in, periodYear)

		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})

	t.Run("год; отнять год от добавочного дня високосного года", func(t *testing.T) {
		in, _ := time.Parse("02.01.2006", "29.02.2020")
		out, _ := time.Parse("02.01.2006", "01.03.2019")

		res, err := fromTime(in, periodYear)

		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})

	t.Run("год; отнять год от следующего дня после добавочного високосного года", func(t *testing.T) {
		in, _ := time.Parse("02.01.2006", "01.03.2020")
		out, _ := time.Parse("02.01.2006", "01.03.2019")

		res, err := fromTime(in, periodYear)

		assert.NoError(t, err)
		assert.Equal(t, out, res)
	})
}
