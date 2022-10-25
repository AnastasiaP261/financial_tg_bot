package db

import (
	"context"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/db/_testdb"
)

func NewTestDB(ctx context.Context, t *testing.T) (testServ *Service, close func()) {
	container := _testdb.NewTestDatabase(ctx, t)

	time.Sleep(5 * time.Second)

	db, err := New(ctx, container)
	assert.NoError(t, err, "NewTestDB")

	err = makeMigrationsUPForTest(ctx, db)
	assert.NoError(t, err)

	return db, func() {
		container.Close(t)
	}
}

func makeMigrationsUPForTest(ctx context.Context, s *Service) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TYPE currency AS ENUM ('RUB', 'EUR', 'USD', 'CNY');

CREATE TABLE users
(
    id   bigint PRIMARY KEY NOT NULL, -- уникальный id юзера
    curr currency           NOT NULL
);

CREATE TABLE categories
(
    id            bigserial PRIMARY KEY NOT NULL, -- уникальный id
    user_id       bigint                NOT NULL,
    category_name text                  NOT NULL
);

CREATE UNIQUE INDEX categories_idx ON categories (user_id, category_name);

CREATE TABLE purchases
(
    id          bigserial PRIMARY KEY NOT NULL, -- уникальный id
    category_id bigint                NOT NULL, -- по этому айди можно узнать и юзера
    sum         numeric               NOT NULL,
    ts          timestamp             NOT NULL DEFAULT NOW(),

    eur_ratio   numeric               NOT NULL,
    usd_ratio   numeric               NOT NULL,
    cny_ratio   numeric               NOT NULL
);

CREATE INDEX purchases_idx ON purchases (category_id, ts);

CREATE TABLE rate
(
    "date"    int PRIMARY KEY NOT NULL, -- уникальный id записи
    eur_ratio numeric         NOT NULL,
    usd_ratio numeric         NOT NULL,
    cny_ratio numeric         NOT NULL
);
`)
	return err
}
