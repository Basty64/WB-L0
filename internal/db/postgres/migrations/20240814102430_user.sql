-- +goose Up
-- +goose StatementBegin

CREATE USER test WITH PASSWORD 'test';
GRANT ALL PRIVILEGES ON DATABASE delivery_service TO test;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

REVOKE ALL PRIVILEGES ON DATABASE delivery_service FROM test;
DROP USER test;

-- +goose StatementEnd
