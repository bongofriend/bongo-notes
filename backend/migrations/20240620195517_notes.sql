-- +goose Up
-- +goose StatementBegin
CREATE TABLE notes(
    id text not null unique,
    notebook_id text not null,
    title text not null,
    path text not null,
    created_at timestamp not null default (strftime('%s','now')), --TODO: Correct data type for timestamps
    updated_at timestamp not null default (strftime('%s','now')),

    FOREIGN KEY (notebook_id) REFERENCES notebooks(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notes;
-- +goose StatementEnd
