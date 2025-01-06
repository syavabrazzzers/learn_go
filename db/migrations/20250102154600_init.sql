-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id serial,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp default NULL,
    email varchar(100) not null unique,
    password varchar not null,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
