-- +goose Up
create table posts (
  id uuid primary key,
  created_at timestamp not null,
  updated_at timestamp not null,
  title varchar(60) not null,
  url varchar(120) unique not null,
  description varchar(300) not null,
  published_at timestamp not null,
  feed_id uuid not null,
  foreign key(feed_id)
    references feeds(id)
    on delete cascade
);

-- +goose Down
drop table posts;
