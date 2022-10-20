package store

import (
	"sync"
)

type flag byte

const (
	flagAdd          = 1
	flagEditCurrency = 2
)

func (s *Service) categoryAccessWrite(c category) {
	s.CategoriesMutex.Lock()

	s.Categories = append(s.Categories, categoryRow{
		RWMutex:  sync.RWMutex{},
		category: c,
	})

	s.CategoriesMutex.Unlock()
}

func (s *Service) purchaseAccessWrite(p purchase) {
	s.PurchasesMutex.Lock()

	s.Purchases = append(s.Purchases, purchaseRow{
		RWMutex:  sync.RWMutex{},
		purchase: p,
	})

	s.PurchasesMutex.Unlock()
}

func (s *Service) usersAccessWrite(u user, flag flag) {
	s.UsersMutex.Lock()
	defer s.UsersMutex.Unlock()

	if flag == flagAdd {
		s.Users = append(s.Users, userRow{
			RWMutex: sync.RWMutex{},
			user:    u,
		})
		return
	}

	if flag == flagEditCurrency {
		for i, row := range s.Users { // nolint:govet
			if row.UserID == u.UserID {
				s.Users[i].Lock()
				s.Users[i].Currency = u.Currency
				s.Users[i].Unlock()
				return
			}
		}
	}
}

// возвращает копию списка категорий (без локов)
func (s *Service) categoryAccessRead() []category {
	s.CategoriesMutex.RLock()

	res := make([]category, len(s.Categories))
	for i := range s.Categories {
		s.Categories[i].RLock()
		res[i] = category{
			UserID:   s.Categories[i].UserID,
			Category: s.Categories[i].Category,
		}
		s.Categories[i].RUnlock()
	}

	s.CategoriesMutex.RUnlock()

	return res
}

// возвращает копию списка трат (без локов)
func (s *Service) purchaseAccessRead() []purchase {
	s.PurchasesMutex.RLock()

	res := make([]purchase, len(s.Purchases))
	for i := range s.Purchases {
		s.Purchases[i].RLock()
		res[i] = purchase{
			UserID:   s.Purchases[i].UserID,
			Sum:      s.Purchases[i].Sum,
			Category: s.Purchases[i].Category,
			Date:     s.Purchases[i].Date,
			USDRatio: s.Purchases[i].USDRatio,
			CNYRatio: s.Purchases[i].CNYRatio,
			EURRatio: s.Purchases[i].EURRatio,
		}
		s.Purchases[i].RUnlock()
	}

	s.PurchasesMutex.RUnlock()

	return res
}

// возвращает копию списка юзеров (без локов)
func (s *Service) usersAccessRead() []user {
	s.UsersMutex.RLock()

	res := make([]user, len(s.Users))
	for i := range s.Users {
		s.Users[i].RLock()
		res[i] = user{
			UserID:   s.Users[i].UserID,
			Currency: s.Users[i].Currency,
		}
		s.Users[i].RUnlock()
	}
	s.UsersMutex.RUnlock()

	return res
}
