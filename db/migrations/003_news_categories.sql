-- +goose Up
create table news_categories(
     news_id integer references news(id) on delete cascade,
     category_id integer,
     unique(news_id, category_id)
);

-- +goose Down
drop table news_categories;