-- +goose Up
-- +goose StatementBegin
create table users (
    id bigserial primary key,
    name VARCHAR(120) not null,
    email VARCHAR(120) not null,
    password VARCHAR(120) not null,
    is_enabled bool default true,

    created_at timestamp without time zone default (now() at time zone 'utc'),
    updated_at timestamp without time zone default (now() at time zone 'utc')
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
Drop table if exists users;

-- +goose StatementEnd
