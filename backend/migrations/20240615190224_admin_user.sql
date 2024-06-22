-- +goose Up
-- +goose StatementBegin
INSERT INTO users(id, username, password) VALUES("896f9ebc-0869-4d52-bf5c-a3071f6a7fef", "admin", "$2a$10$k3y2BdiILp0jafVOVK9i5O1zoDX9QTHuxYvsjIcTAfJhDU6N.Srpa");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users Where username = "admin";
-- +goose StatementEnd
