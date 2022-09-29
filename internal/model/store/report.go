package store

import (
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
	"time"
)

func (s *Service) GetReport(fromDate time.Time, userID int64) ([]model.ReportItem, error) {
	purchases := make(map[string]float64, 0)
	for _, p := range s.Purchases {
		if p.UserID == userID {
			if p.Date.After(fromDate) || p.Date.Equal(fromDate) {
				purchases[p.Category] += p.Sum
			}
		}
	}

	res := make([]model.ReportItem, 0, len(purchases))
	for k, v := range purchases {
		res = append(res, model.ReportItem{
			PurchaseCategory: k,
			Summa:            v,
		})
	}

	return res, nil
}
