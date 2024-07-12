-- +goose Up
create table feed_user (
  id uuid primary key,
  feed_id uuid not null,
  user_id uuid not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  foreign key(feed_id) references feeds(id),
  foreign key(user_id) references users(id)
);

-- +goose Down
drop table feed_user;
