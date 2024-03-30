-- +goose Up
create table users(
      id uuid primary key,
      name varchar not null,
      email varchar unique not null,
      password varchar not null,
      updated_at timestamptz not null default now()
);

-- +goose Down
drop table users;