package store

import (
	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	"testing"
)

func Test_CategoryExist(t *testing.T) {
	s := New()

	s.Categories = []category{
		{UserID: 123, Category: "some category"},
	}

	res := s.CategoryExist(model.CategoryExistReq{
		UserID:   123,
		Category: "some category",
	})

	assert.True(t, res)
}
