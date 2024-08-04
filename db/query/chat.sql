-- name: InsertMessage :one
INSERT INTO message (from_user_id, to_user_id, is_sent, content)
VALUES ($1, $2, $3, $4)
RETURNING *;



-- name: GetMessages :many
SELECT * FROM message
WHERE from_user_id = $1
AND to_user_id = $2;
