-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
id UUID NOT NULL,
passport_series INTEGER NOT NULL,
passport_number INTEGER NOT NULL,
surname TEXT NOT NULL,
name TEXT NOT NULL,
patronymic TEXT NOT NULL,
address TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
