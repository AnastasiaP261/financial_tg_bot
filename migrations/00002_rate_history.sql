-- +goose Up

-- здесь автоматически создастся уникальный индекс для id
CREATE TABLE rate
(
    "date"    int PRIMARY KEY NOT NULL, -- уникальный id записи
    eur_ratio numeric         NOT NULL,
    usd_ratio numeric         NOT NULL,
    cny_ratio numeric         NOT NULL
);

-- +goose Down

-- индексы удалятся автоматически при удалении таблиц
DROP TABLE rate;
