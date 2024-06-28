-- +goose Up
-- +goose StatementBegin
CREATE TABLE note_diffs(
    id text not null,
    note_id text not null,
    created_at timestamp not null default (strftime('%s','now')),

    FOREIGN KEY(note_id) REFERENCES notes(id)
);
CREATE UNIQUE INDEX idx_diffs_created_at on note_diffs(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_diffs_created_at;
DROP TABLE note_diffs;
-- +goose StatementEnd
