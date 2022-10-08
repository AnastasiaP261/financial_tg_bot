package store

type flag byte

const (
	flagAdd          = 1
	flagEditCurrency = 2
)

func (s *Service) categoryAccessWrite(c category) {
	s.CategoriesMutex.Lock()

	s.Categories = append(s.Categories, c)

	s.CategoriesMutex.Unlock()
}

func (s *Service) purchaseAccessWrite(p purchase) {
	s.PurchasesMutex.Lock()

	s.Purchases = append(s.Purchases, p)

	s.PurchasesMutex.Unlock()
}

func (s *Service) usersAccessWrite(u user, flag flag) {
	s.UsersMutex.Lock()
	defer s.UsersMutex.Unlock()

	if flag == flagAdd {
		s.Users = append(s.Users, u)
		return
	}

	if flag == flagEditCurrency {
		for i, row := range s.Users {
			if row.UserID == u.UserID {
				s.Users[i].Currency = u.Currency
				return
			}
		}
	}
}
