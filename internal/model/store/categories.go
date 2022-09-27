package store

import model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"

func (s *Service) CategoryExist(req model.CategoryExistReq) bool {
	for _, item := range s.Categories {
		if item == category(req) {
			return true
		}
	}
	return false
}
