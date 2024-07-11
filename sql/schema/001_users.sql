-- +goose up 
create table users(
  id uuid primary key, 
  created_at timestamp not null,
  updated_at timestamp not null,
  name varchar(40) not null
);

-- +goose down
