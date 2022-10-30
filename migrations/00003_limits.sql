-- +goose Up

ALTER TABLE users ADD COLUMN month_limit numeric NOT NULL DEFAULT -1; -- -1 значит не установлен лимит

-- +goose Down

ALTER TABLE users DROP COLUMN month_limit;

