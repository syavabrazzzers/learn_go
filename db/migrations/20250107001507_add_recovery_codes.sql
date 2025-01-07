-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_recovery_codes (
    user_id integer,
    codes jsonb,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    PRIMARY KEY (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_recovery_codes;
-- +goose StatementEnd
