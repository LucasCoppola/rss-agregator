-- name: CreateFeed :one
insert into feeds(id, created_at, updated_at, name, url, user_id)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetFeeds :many
select * from feeds;

-- name: GetNextFeedsToFetch :many
select *
from feeds
order by last_fetched_at is not null, last_fetched_at
limit $1;

-- name: MarkFeedFetched :exec
update feeds 
set last_fetched_at = $1, updated_at = $2 
where id = $3;
