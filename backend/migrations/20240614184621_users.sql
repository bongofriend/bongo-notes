-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    username text not null,
    password text not null,
    created_at timestamp not null default (strftime('%s','now')), --TODO: Correct data type for timestamps
    updated_at timestamp not null default (strftime('%s','now'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
