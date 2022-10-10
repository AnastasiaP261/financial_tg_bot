package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func Test_ChangeCurrency(t *testing.T) {
	s := New()
	s.Users = []user{
		{UserID: 123, Currency: RUB},
	}

	err := s.ChangeCurrency(123, purchases.USD)

	assert.NoError(t, err)
	assert.Equal(t, []user{{
		UserID:   123,
		Currency: USD,
	}}, s.Users)
}

func Test_addUser(t *testing.T) {
	s := New()
	s.Users = []user{}

	err := s.addUser(123)

	assert.NoError(t, err)
	assert.Equal(t, []user{{
		UserID:   123,
		Currency: RUB,
	}}, s.Users)
}

func Test_getUserInfo(t *testing.T) {
	s := New()
	s.Users = []user{
		{UserID: 123, Currency: RUB},
	}

	res, err := s.getUserInfo(123)

	assert.NoError(t, err)
	assert.Equal(t, user{
		UserID:   123,
		Currency: RUB,
	}, res)
}

func Test_userExist(t *testing.T) {
	t.Run("юзер есть", func(t *testing.T) {
		s := New()
		s.Users = []user{
			{UserID: 123, Currency: RUB},
		}

		res, err := s.userExist(123)

		assert.NoError(t, err)
		assert.True(t, res)
	})

	t.Run("юзера нет", func(t *testing.T) {
		s := New()
		s.Users = []user{}

		res, err := s.userExist(123)

		assert.NoError(t, err)
		assert.False(t, res)
	})
}
