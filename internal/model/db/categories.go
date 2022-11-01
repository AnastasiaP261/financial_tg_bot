package db

import (
	"context"
	"database/sql"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type category struct { // nolint:unused
	ID       uint64 `db:"id"`
	Category string `db:"category_name"`
}

// GetCategoryID получить id категории
func (s *Service) GetCategoryID(ctx context.Context, categoryName string) (uint64, error) {
	if categoryName == "" {
		return 0, errors.New("category is empty")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tblCategoriesColID).
		From(tblCategories).
		Where(sq.Eq{
			tblCategoriesColCategoryName: categoryName,
		}).
		ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "query creating error")
	}

	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, errors.Wrap(err, "db.QueryRowContext")
	}
	var id uint64
	if err = read(rows, &id); err != nil {
		return 0, errors.Wrap(err, "read")
	}

	return id, nil
}

func (s *Service) AddCategory(ctx context.Context, categoryName string) error {
	if categoryName == "" {
		return errors.New("category is empty")
	}

	categoryID, err := s.GetCategoryID(ctx, categoryName)
	if err != nil {
		return errors.Wrap(err, "GetCategoryID")
	}
	if categoryID != 0 {
		return ErrCategoryAlreadyExists
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tblCategories).
		Columns(tblCategoriesColCategoryName).
		Values(categoryName).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args...); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

// GetAllCategories получить все категории
func (s *Service) GetAllCategories(ctx context.Context) ([]model.CategoryRow, error) {
	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tblCategoriesColID, tblCategoriesColCategoryName).
		From(tblCategories).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query creating error")
	}

	var categories []category
	if err = s.db.SelectContext(ctx, &categories, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "db.QueryRowContext")
	}

	res := make([]model.CategoryRow, len(categories))
	for i := range categories {
		res[i] = model.CategoryRow(categories[i])
	}

	return res, nil
}
