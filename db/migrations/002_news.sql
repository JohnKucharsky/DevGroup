-- +goose Up
create table news(
      id serial primary key,
      title varchar not null,
      content varchar unique not null,
      updated_at timestamptz not null default now()
);

-- +goose Down
drop table news;