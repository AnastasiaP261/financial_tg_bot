package store

import (
	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	"testing"
	"time"
)

func Test_AddPurchase(t *testing.T) {
	s := New()

	date, _ := time.Parse("02.01.2006", "01.01.2000")
	err := s.AddPurchase(model.AddPurchaseReq{
		UserID:   123,
		Sum:      100.50,
		Category: "some category",
		Date:     date,
	})

	assert.NoError(t, err)
	assert.Equal(t,
		s.Purchases[0],
		purchase{
			UserID:   123,
			Sum:      100.50,
			Category: "some category",
			Date:     date,
		},
	)
}
