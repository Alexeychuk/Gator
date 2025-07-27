-- name: CreateFeedFollow :one
WITH inserted_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    inserted_follow.*,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_follow
JOIN users u ON inserted_follow.user_id = u.id
JOIN feeds f ON inserted_follow.feed_id = f.id;


-- name: GetFeedFollows :many
SELECT * FROM feed_follows;

-- name: GetFeedFollowsForUser :many
SELECT 
    ff.*, 
    u.name as user_name, 
    f.name as feed_name 
FROM feed_follows ff
JOIN users u ON ff.user_id = u.id
JOIN feeds f ON ff.feed_id = f.id
WHERE ff.user_id = $1;

-- name: DeleteFeedFollowByUserIdAndUrl :exec
DELETE FROM feed_follows 
WHERE feed_follows.user_id = $1
AND feed_id = (
    SELECT id FROM feeds 
    WHERE feeds.url = $2
);