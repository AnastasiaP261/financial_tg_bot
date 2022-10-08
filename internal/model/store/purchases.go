package store

import (
	"time"

	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func (s *Service) AddPurchase(req model.AddPurchaseReq) error {
	if err := s.UserCreateIfNotExist(req.UserID); err != nil {
		return errors.Wrap(err, "UserCreateIfNotExist")
	}

	if req.UserID == 0 {
		return errors.New("user is empty")
	}
	if req.Sum == 0 {
		return errors.New("sum is empty")
	}
	{
		nilTime := time.Time{}
		if req.Date == nilTime {
			return errors.New("date is empty")
		}
	}

	s.purchaseAccessWrite(purchase(req))

	return nil
}

// GetUserPurchasesFromDate получить все траты пользователя
func (s *Service) GetUserPurchasesFromDate(fromDate time.Time, userID int64) ([]model.Purchase, error) {
	if err := s.UserCreateIfNotExist(userID); err != nil {
		return nil, errors.Wrap(err, "UserCreateIfNotExist")
	}

	purchases := make([]model.Purchase, 0)
	for _, p := range s.Purchases {
		if p.UserID == userID {
			if p.Date.After(fromDate) || p.Date.Equal(fromDate) {
				purchases = append(purchases, model.Purchase{
					PurchaseCategory: p.Category,
					Summa:            p.Sum,
					RateToRUB: model.RateToRUB{
						USD: p.USDRatio,
						CNY: p.CNYRatio,
						EUR: p.EURRatio,
					},
				})

			}
		}
	}

	return purchases, nil
}
