// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH
    inserted_feed_follows AS (
        INSERT INTO
            feed_follows (
                id,
                created_at,
                updated_at,
                user_id,
                feed_id
            )
        VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at, user_id, feed_id
    )
SELECT inserted_feed_follows.id, inserted_feed_follows.created_at, inserted_feed_follows.updated_at, inserted_feed_follows.user_id, inserted_feed_follows.feed_id, feeds.name AS feed_name, users.name AS user_name
FROM
    inserted_feed_follows
    INNER JOIN feeds ON inserted_feed_follows.feed_id = feeds.id
    INNER JOIN users ON inserted_feed_follows.user_id = users.id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	UserName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const createFeeds = `-- name: CreateFeeds :one
INSERT INTO
    feeds (
        id,
        created_at,
        updated_at,
        name,
        url,
        user_id
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at, name, url, user_id
`

type CreateFeedsParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) CreateFeeds(ctx context.Context, arg CreateFeedsParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeeds,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const deleteFollowByUserAndUrl = `-- name: DeleteFollowByUserAndUrl :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2
`

type DeleteFollowByUserAndUrlParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) DeleteFollowByUserAndUrl(ctx context.Context, arg DeleteFollowByUserAndUrlParams) error {
	_, err := q.db.ExecContext(ctx, deleteFollowByUserAndUrl, arg.UserID, arg.FeedID)
	return err
}

const getFeedByUrl = `-- name: GetFeedByUrl :one
SELECT id, created_at, updated_at, name, url, user_id FROM feeds WHERE url = $1
`

func (q *Queries) GetFeedByUrl(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUrl, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const getFeedFollowsByUsers = `-- name: GetFeedFollowsByUsers :many
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, feed_follows.feed_id, users.name as user_name, feeds.name AS feed_name
FROM
    feed_follows
    INNER JOIN users ON users.id = feed_follows.user_id
    INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
    users.name = $1
`

type GetFeedFollowsByUsersRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	UserName  string
	FeedName  string
}

func (q *Queries) GetFeedFollowsByUsers(ctx context.Context, name string) ([]GetFeedFollowsByUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsByUsers, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsByUsersRow
	for rows.Next() {
		var i GetFeedFollowsByUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.UserName,
			&i.FeedName,
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

const getFeeds = `-- name: GetFeeds :many
SELECT
    feeds.name AS feed_name,
    feeds.url AS feed_url,
    users.name AS user_name
FROM feeds
    LEFT JOIN users ON feeds.user_id = users.id
`

type GetFeedsRow struct {
	FeedName string
	FeedUrl  string
	UserName sql.NullString
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(&i.FeedName, &i.FeedUrl, &i.UserName); err != nil {
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
