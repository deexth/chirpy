-- name: CreateChirp :one
INSERT INTO chirps (
    id,
    body,
    user_id
) VALUES ( $1, $2, $3 ) RETURNING *;

-- name: GetChirps :many
SELECT *
    FROM chirps
    ORDER BY created_at ASC
    LIMIT $1;

-- name: GetChirpsByAuthorID :many
SELECT *
    FROM chirps
    WHERE user_id = $1
    ORDER BY created_at ASC
    LIMIT $2;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirp :execrows
DELETE FROM chirps
    WHERE id = $1 AND user_id = $2;
