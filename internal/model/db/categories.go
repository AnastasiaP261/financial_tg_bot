package db

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

type category struct {
	ID       uint64 `db:"id"`
	UserID   uint64 `db:"user_id"`
	Category string `db:"category_name"`
}

// GetCategoryID получить id категории
func (s *Service) GetCategoryID(ctx context.Context, req model.CategoryRow) (uint64, error) {
	fmt.Println("### GetCategoryID")
	if req.UserID == 0 {
		return 0, errors.New("userID is empty")
	}
	if req.Category == "" {
		return 0, errors.New("category is empty")
	}

	if err := s.UserCreateIfNotExist(ctx, req.UserID); err != nil {
		return 0, errors.Wrap(err, "UserCreateIfNotExist")
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
		return 0, errors.Wrap(err, "query creating error")
	}

	fmt.Println("### q", q, args)

	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, errors.Wrap(err, "db.QueryRowContext")
	}
	var id uint64
	read(rows, &id)

	fmt.Println("### id", id)

	return id, nil
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

	categoryID, err := s.GetCategoryID(ctx, req)
	if err != nil {
		return errors.Wrap(err, "GetCategoryID")
	}
	if categoryID != 0 {
		return ErrCategoryAlreadyExists
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tblCategories).
		Columns(tblCategoriesColCategoryName, tblCategoriesColUserID).
		Values(req.Category, req.UserID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args...); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}
