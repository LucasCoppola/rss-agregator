// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: follow_feed.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const followFeed = `-- name: FollowFeed :one
insert into feed_user (id, feed_id, user_id, created_at, updated_at)
values ($1, $2, $3, $4, $5)
returning id, feed_id, user_id, created_at, updated_at
`

type FollowFeedParams struct {
	ID        uuid.UUID
	FeedID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) FollowFeed(ctx context.Context, arg FollowFeedParams) (FeedUser, error) {
	row := q.db.QueryRowContext(ctx, followFeed,
		arg.ID,
		arg.FeedID,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i FeedUser
	err := row.Scan(
		&i.ID,
		&i.FeedID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFollowedFeeds = `-- name: GetFollowedFeeds :many
select id, feed_id, user_id, created_at, updated_at from feed_user 
where user_id = $1
`

func (q *Queries) GetFollowedFeeds(ctx context.Context, userID uuid.UUID) ([]FeedUser, error) {
	rows, err := q.db.QueryContext(ctx, getFollowedFeeds, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedUser
	for rows.Next() {
		var i FeedUser
		if err := rows.Scan(
			&i.ID,
			&i.FeedID,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unfollowFeed = `-- name: UnfollowFeed :exec
delete from feed_user where id = $1 and user_id = $2
`

type UnfollowFeedParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) UnfollowFeed(ctx context.Context, arg UnfollowFeedParams) error {
	_, err := q.db.ExecContext(ctx, unfollowFeed, arg.ID, arg.UserID)
	return err
}
