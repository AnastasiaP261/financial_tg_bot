package store

import (
	"github.com/pkg/errors"
	"time"

	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func (s *Service) GetReport(fromDate time.Time, userID int64) ([]model.ReportItem, error) {
	if err := s.UserCreateIfNotExist(userID); err != nil {
		return nil, errors.Wrap(err, "UserCreateIfNotExist")
	}

	info, err := s.getUserInfo(userID)
	if err != nil {
		return nil, errors.Wrap(err, "getUserInfo")
	}
	cy := info.Currency

	purchases := make(map[string]float64, 0)
	for _, p := range s.Purchases {
		if p.UserID == userID {
			var ratio float64
			switch cy {
			case RUB:
				ratio = 1
			case EUR:
				ratio = p.EURRatio
			case USD:
				ratio = p.USDRatio
			case CNY:
				ratio = p.CNYRatio
			}

			if p.Date.After(fromDate) || p.Date.Equal(fromDate) {
				purchases[p.Category] += p.Sum * ratio
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
