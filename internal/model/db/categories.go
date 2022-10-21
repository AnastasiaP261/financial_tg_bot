package db

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

// CategoryExist проверка, что категория существует
func (s *Service) CategoryExist(ctx context.Context, req model.CategoryRow) (bool, error) {
	fmt.Println("### CategoryExist")
	if req.UserID == 0 {
		return false, errors.New("userID is empty")
	}
	if req.Category == "" {
		return false, errors.New("category is empty")
	}

	if err := s.UserCreateIfNotExist(ctx, req.UserID); err != nil {
		return false, errors.Wrap(err, "UserCreateIfNotExist")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tblCategoriesColID).
		From(tblCategories).
		Where(sq.Eq{
			tblCategoriesColUserID:       req.UserID,
			tblCategoriesColCategoryName: req.Category,
		}).
		ToSql()
	if err != nil {
		return false, errors.Wrap(err, "query creating error")
	}

	var id uint64
	if err = s.db.QueryRowContext(ctx, q, args).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, "db.QueryRowContext")
	}

	return false, nil
}

func (s *Service) AddCategory(ctx context.Context, req model.CategoryRow) error {
	fmt.Println("### AddCategory")
	if req.UserID == 0 {
		return errors.New("userID is empty")
	}
	if req.Category == "" {
		return errors.New("category is empty")
	}

	if err := s.UserCreateIfNotExist(ctx, req.UserID); err != nil {
		return errors.Wrap(err, "UserCreateIfNotExist")
	}

	exist, err := s.CategoryExist(ctx, req)
	if err != nil {
		return errors.Wrap(err, "CategoryExist")
	}
	if exist {
		return errors.New("category is already exist")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tblCategories).
		Columns(tblCategoriesColCategoryName, tblCategoriesColUserID).
		Values(req.Category, req.UserID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}
