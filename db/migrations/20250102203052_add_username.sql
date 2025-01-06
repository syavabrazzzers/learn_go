-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN username varchar null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE DROP COLUMN username;
-- +goose StatementEnd
