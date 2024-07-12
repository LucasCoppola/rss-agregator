-- name: FollowFeed :one
insert into feed_user (id, feed_id, user_id, created_at, updated_at)
values ($1, $2, $3, $4, $5)
returning *;

-- name: UnfollowFeed :exec
delete from feed_user where id = $1 and user_id = $2;
