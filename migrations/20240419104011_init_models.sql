-- +goose Up
-- +goose StatementBegin

CREATE TABLE url
(
    id BIGSERIAL PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    URL TEXT NOT NULL
);

CREATE INDEX url_alias ON url(alias);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX url_alias;
DROP TABLE url;

-- +goose StatementEnd
