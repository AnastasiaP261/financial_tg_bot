-- +goose Up

CREATE TYPE currency AS ENUM ('RUB', 'EUR', 'USD', 'CNY');

-- здесь автоматически создастся уникальный индекс для id
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

-- будем искать либо по id (но на него индекс создается автоматически),
-- либо по паре юзер+категория, b-tree подходит
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

-- будем искать по категории и сравнивать с заданной датой, b-tree подходит
CREATE INDEX purchases_idx ON purchases (category_id, ts);

-- +goose Down

-- индексы удалятся автоматически при удалении таблиц
DROP TABLE purchases;
DROP TABLE categories;
DROP TABLE users;
DROP TYPE currency;
