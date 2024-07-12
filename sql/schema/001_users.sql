-- +goose up 
create table users(
  id uuid primary key, 
  created_at timestamp not null,
  updated_at timestamp not null,
  name varchar(40) not null,
  api_key varchar(64) unique not null default(
    encode(sha256(random()::text::bytea), 'hex')
  )
);

-- +goose down
