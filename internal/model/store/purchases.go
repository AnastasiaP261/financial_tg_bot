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
