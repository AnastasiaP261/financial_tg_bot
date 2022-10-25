package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	model "gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/purchases"
)

// rate курс валют к RUB
type rate struct {
	USD float64 `db:"usd_ratio"`
	EUR float64 `db:"eur_ratio"`
	CNY float64 `db:"cny_ratio"`
}

// GetRate проверяет есть ли в базе курсы валют за указанную дату. если есть, возвращает true и их значения
func (s *Service) GetRate(ctx context.Context, y, m, d int) (bool, model.RateToRUB, error) {
	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(tblRateColCNYRatio, tblRateColUSDRatio, tblRateColEURRatio).
		From(tblRate).
		Where(sq.Eq{tblRateColDate: dateToKey(y, m, d)}).
		ToSql()

	if err != nil {
		return false, model.RateToRUB{}, errors.Wrap(err, "query creating error")
	}

	rows, err := s.db.QueryxContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, model.RateToRUB{}, nil
		}
		return false, model.RateToRUB{}, errors.Wrap(err, "db.QueryxContext")
	}

	var r rate
	if err = readX(rows, &r); err != nil {
		return false, model.RateToRUB{}, errors.Wrap(err, "readX")
	}

	emptyRates := rate{}
	if r == emptyRates {
		return false, model.RateToRUB{}, nil
	}

	return true, model.RateToRUB(r), nil
}

func (s *Service) AddRate(ctx context.Context, y, m, d int, rates model.RateToRUB) error {
	q, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(tblRate).
		Columns(tblRateColDate, tblRateColEURRatio, tblRateColUSDRatio, tblRateColCNYRatio).
		Values(dateToKey(y, m, d), rates.EUR, rates.USD, rates.CNY).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "query creating error")
	}

	if _, err = s.db.ExecContext(ctx, q, args...); err != nil {
		return errors.Wrap(err, "db.ExecContext")
	}

	return nil
}

func dateToKey(y, m, d int) int {
	return y*10000 + m*100 + d
}
