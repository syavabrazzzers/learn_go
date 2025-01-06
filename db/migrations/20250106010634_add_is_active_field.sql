-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP COLUMN is_active FROM users;
-- +goose StatementEnd
