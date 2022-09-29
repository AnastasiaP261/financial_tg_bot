package store

import (
	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

func (s *Service) CategoryExist(req model.CategoryRow) (bool, error) {
	if req.UserID == 0 {
		return false, errors.New("userID is empty")
	}
	if req.Category == "" {
		return false, errors.New("category is empty")
	}

	for _, item := range s.Categories {
		if item == category(req) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) AddCategory(req model.CategoryRow) error {
	if req.UserID == 0 {
		return errors.New("userID is empty")
	}
	if req.Category == "" {
		return errors.New("category is empty")
	}

	exist, err := s.CategoryExist(req)
	if err != nil {
		return errors.Wrap(err, "CategoryExist")
	}
	if exist {
		return errors.New("category is already exist")
	}

	s.categoryAccessWrite(category{
		UserID:   req.UserID,
		Category: req.Category,
	})

	return nil
}
