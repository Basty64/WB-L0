-- +goose Up
-- +goose StatementBegin

ALTER USER admin WITH PASSWORD 'admin';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS deliveries CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS items CASCADE;

-- +goose StatementEnd
