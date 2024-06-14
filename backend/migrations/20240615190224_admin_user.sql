-- +goose Up
-- +goose StatementBegin
INSERT INTO users(username, password) VALUES("admin", "$2a$10$k3y2BdiILp0jafVOVK9i5O1zoDX9QTHuxYvsjIcTAfJhDU6N.Srpa");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users Where username = "admin";
-- +goose StatementEnd
