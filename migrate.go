package migrate

import (
	"database/sql"
	"embed"
	"io/fs"

	goose "github.com/pressly/goose/v3"
)

// 1) паттерн пути к файлу или директории не должен начинаться с / и иметь в себе . или ..
// 2) паттерн может принимать только файлы или директории внутри модуля, но не во вне
// Т.е. для чтения миграций директива должна читаться в файле в корневом каталоге
// https://nikgalushko.com/post/go_1.16_embed/

//go:embed migrations/*.sql
var Migrations embed.FS

func Migrate(dsn string, path fs.FS) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	goose.SetBaseFS(path)
	return goose.Up(db, "migrations")
}
