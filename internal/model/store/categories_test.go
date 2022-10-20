package store

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func Test_CategoryExist(t *testing.T) {
	s := New()

	s.Users = []userRow{
		{sync.RWMutex{}, user{UserID: 123, Currency: RUB}},
	}
	s.Categories = []categoryRow{
		{sync.RWMutex{}, category{UserID: 123, Category: "some category"}},
	}

	res, err := s.CategoryExist(model.CategoryRow{
		UserID:   123,
		Category: "some category",
	})

	assert.NoError(t, err)
	assert.True(t, res)
}
