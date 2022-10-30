package _testdb

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	testcon "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot"
)

const timeout = time.Minute

type TestDatabase struct {
	instance testcon.Container
	uri      string
}

// NewTestDatabase запускает в контейнере базу для проведения интеграционных тестов
// и возвращает ее инстанс
//
// usage:
// 	testDB := _testdb.NewTestDatabase(сtx, t)
// 	defer testDB.Close(t)
// 	println(testDB.DBUri())
func NewTestDatabase(ctx context.Context, t *testing.T) *TestDatabase {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req := testcon.ContainerRequest{
		Image:        "postgres:12",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgres, err := testcon.GenericContainer(ctx, testcon.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err, "testcon.GenericContainer")

	tdb := &TestDatabase{
		instance: postgres,
	}
	tdb.uri = fmt.Sprintf("postgres://postgres:postgres@127.0.0.1:%d/postgres?sslmode=disable", tdb.port(t))

	// run migrations
	err = migrate.Migrate(tdb.uri, migrate.Migrations)

	return tdb
}

func (db *TestDatabase) port(t *testing.T) int {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p, err := db.instance.MappedPort(ctx, "5432")
	assert.NoError(t, err, "instance.MappedPort")
	return p.Int()
}

func (db *TestDatabase) DBUri() string {
	return db.uri
}

func (db *TestDatabase) Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	assert.NoError(t, db.instance.Terminate(ctx), "db.instance.Terminate(ctx)")
}

func (db *TestDatabase) Exec(ctx context.Context, cmd []string) (int, io.Reader, error) {
	return db.instance.Exec(ctx, cmd)
}
