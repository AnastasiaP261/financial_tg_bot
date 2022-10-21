package db

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"log"
	"time"

	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

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
		Values(req.Category, req.Sum, req.EURRatio, req.USDRatio, req.CNYRatio)

	nullTime := time.Time{}
	if req.Date == nullTime {
		query.Columns(tblPurchasesColTimestamp).Values(req.Date)
	}

	q, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

type purchase struct {
	Sum      float64 `db:"sum"` // сумма траты в рублях
	Category string  `db:"category"`

	// коэффициенты валют на момент совершения траты
	USDRatio float64 `db:"usd_ratio"`
	CNYRatio float64 `db:"cny_ratio"`
	EURRatio float64 `db:"eur_ratio"`
}

// GetUserPurchasesFromDate получить все траты пользователя
func (s *Service) GetUserPurchasesFromDate(ctx context.Context, fromDate time.Time, userID int64) ([]model.Purchase, error) {
	if err := s.UserCreateIfNotExist(ctx, userID); err != nil {
		return nil, errors.Wrap(err, "UserCreateIfNotExist")
	}

	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tblPurchasesColSum, tblPurchasesColCategoryID,
			tblPurchasesColUSDRatio, tblPurchasesColCNYRatio, tblPurchasesColEURRatio).
		From(tblPurchases).
		Where(sq.And{
			sq.Eq{tblUsersColID: userID},
			sq.Expr("$1 <= $2", tblPurchasesColTimestamp, fromDate),
		}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "query creating error")
	}

	rows, err := s.db.QueryContext(ctx, q, args)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	purchases := make([]model.Purchase, 0)
	for rows.Next() {
		var p purchase
		if err := rows.Scan(&p); err != nil {
			log.Fatal(err)
		}

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

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return purchases, nil
}
