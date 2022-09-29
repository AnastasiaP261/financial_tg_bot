package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func Test_CategoryExist(t *testing.T) {
	s := New()

	s.Categories = []category{
		{UserID: 123, Category: "some category"},
	}

	res, err := s.CategoryExist(model.CategoryRow{
		UserID:   123,
		Category: "some category",
	})

	assert.NoError(t, err)
	assert.True(t, res)
}
