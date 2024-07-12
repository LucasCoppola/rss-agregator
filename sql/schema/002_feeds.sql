-- +goose Up
create table feeds (
  id uuid primary key,
  created_at timestamp not null,
  updated_at timestamp not null,
  name varchar(60) not null,
  url varchar(120) unique not null,
  user_id uuid not null,
  constraint fk_user
    foreign key(user_id)
    references users(id)
    on delete cascade
);

-- +goose Down
drop table feeds;
