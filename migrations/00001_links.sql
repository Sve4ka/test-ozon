-- +goose Up
CREATE TABLE IF NOT EXISTS links
(
    id           SERIAL PRIMARY KEY,
    original_link VARCHAR     NOT NULL UNIQUE,
    short_code VARCHAR(10) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS links;

