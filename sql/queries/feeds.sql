-- name: CreateFeeds :one
INSERT INTO
    feeds (
        id,
        created_at,
        updated_at,
        name,
        url,
        user_id
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetFeeds :many
SELECT
    feeds.name AS feed_name,
    feeds.url AS feed_url,
    users.name AS user_name
FROM feeds
    LEFT JOIN users ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one
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
        VALUES ($1, $2, $3, $4, $5) RETURNING *
    )
SELECT inserted_feed_follows.*, feeds.name AS feed_name, users.name AS user_name
FROM
    inserted_feed_follows
    INNER JOIN feeds ON inserted_feed_follows.feed_id = feeds.id
    INNER JOIN users ON inserted_feed_follows.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedFollowsByUsers :many
SELECT feed_follows.*, users.name as user_name, feeds.name AS feed_name
FROM
    feed_follows
    INNER JOIN users ON users.id = feed_follows.user_id
    INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
    users.name = $1;

-- name: DeleteFollowByUserAndUrl :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;