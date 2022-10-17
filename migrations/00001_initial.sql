-- +goose Up

CREATE TABLE test
(
    id bigserial PRIMARY KEY -- уникальный id
);

-- +goose Down

DROP TABLE test;
