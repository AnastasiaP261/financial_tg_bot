-- +goose Up

ALTER TABLE users ADD COLUMN month_limit numeric NOT NULL DEFAULT -1; -- -1 значит не установлен лимит
ALTER TABLE users ADD COLUMN category_ids bigint[] NOT NULL DEFAULT ARRAY[1]; -- 1 для дефолтной категории
ALTER TABLE categories DROP COLUMN user_id;
ALTER TABLE purchases ADD COLUMN user_id bigint NOT NULL;

ALTER TABLE users ALTER COLUMN curr SET DEFAULT 'RUB';

DROP INDEX IF EXISTS categories_idx;
CREATE UNIQUE INDEX categories_idx ON categories (id, category_name);

INSERT INTO categories (category_name) VALUES ('Не заданная категория');

-- +goose Down

TRUNCATE TABLE categories;

ALTER TABLE users DROP COLUMN month_limit;
ALTER TABLE users DROP COLUMN category_ids;
ALTER TABLE categories ADD COLUMN user_id bigint NOT NULL;
ALTER TABLE purchases DROP COLUMN user_id;

DROP INDEX IF EXISTS categories_idx;
CREATE UNIQUE INDEX categories_idx ON categories (user_id, category_name);

