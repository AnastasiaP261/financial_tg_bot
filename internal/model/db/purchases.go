package db

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

type purchase struct {
	Sum          float64        `db:"sum"` // сумма траты в рублях
	CategoryName sql.NullString `db:"category_name"`

	// коэффициенты валют на момент совершения траты
	USDRatio float64 `db:"usd_ratio"`
	CNYRatio float64 `db:"cny_ratio"`
	EURRatio float64 `db:"eur_ratio"`
}

func (s *Service) AddPurchase(ctx context.Context, req model.AddPurchaseReq) error {
	if err := s.UserCreateIfNotExist(ctx, req.UserID); err != nil {
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

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tblPurchases).
		Columns(tblPurchasesColCategoryID, tblPurchasesColSum, tblPurchasesColEURRatio,
			tblPurchasesColUSDRatio, tblPurchasesColCNYRatio).
		Values(req.CategoryID, req.Sum, req.EURRatio, req.USDRatio, req.CNYRatio)

	nullTime := time.Time{}
	if req.Date == nullTime {
		query.Columns(tblPurchasesColTimestamp).Values(sq.Expr("now()"))
	}

	q, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	_, err = s.db.ExecContext(ctx, q, args...)
	if err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

// GetUserPurchasesFromDate получить все траты пользователя
func (s *Service) GetUserPurchasesFromDate(ctx context.Context, fromDate time.Time, userID int64) ([]model.Purchase, error) {
	if err := s.UserCreateIfNotExist(ctx, userID); err != nil {
		return nil, errors.Wrap(err, "UserCreateIfNotExist")
	}

	q, args, err := sq.Expr(`SELECT "sum", category_name, usd_ratio, cny_ratio, eur_ratio 
							FROM purchases 
							LEFT JOIN (
								SELECT id, category_name 
								FROM categories 
								WHERE user_id=$1 
							) AS user_categories ON (purchases.category_id=user_categories.id) 
							WHERE ts >= $2;`, userID, fromDate).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query creating error")
	}

	var rows []purchase
	if err = s.db.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, errors.Wrap(err, "db.SelectContext")
	}

	purchases := make([]model.Purchase, 0)
	for _, p := range rows {
		purchases = append(purchases, model.Purchase{
			PurchaseCategory: p.CategoryName.String,
			Summa:            p.Sum,
			RateToRUB: model.RateToRUB{
				USD: p.USDRatio,
				CNY: p.CNYRatio,
				EUR: p.EURRatio,
			},
		})
	}

	return purchases, nil
}
