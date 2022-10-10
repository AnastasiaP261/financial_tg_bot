package store

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
