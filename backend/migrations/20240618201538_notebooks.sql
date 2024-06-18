-- +goose Up
-- +goose StatementBegin
CREATE TABLE notebooks(
    creater_id int not null,
    title text not null,
    description text not null,
    created_at timestamp not null default (strftime('%s','now')), --TODO: Correct data type for timestamps
    updated_at timestamp not null default (strftime('%s','now')),
    FOREIGN KEY(creater_id) REFERENCES users(rowid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notebooks;
-- +goose StatementEnd
