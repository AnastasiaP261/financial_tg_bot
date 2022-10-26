package db

import (
	"context"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/model/db/_testdb"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func NewTestDB(ctx context.Context, t *testing.T) (testServ *Service, close func()) {
	container := _testdb.NewTestDatabase(ctx, t)

	time.Sleep(5 * time.Second)

	db, err := New(ctx, container)
	assert.NoError(t, err, "NewTestDB")

	return db, func() {
		container.Close(t)
	}
}
