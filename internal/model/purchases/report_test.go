//go:build test_all || unit_test

package purchases

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
