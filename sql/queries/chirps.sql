-- name: CreateChirp :one
INSERT INTO chirps (
    id,
    body,
    user_id
) VALUES ( $1, $2, $3 ) RETURNING *;

-- name: GetChirps :many
SELECT *
    FROM chirps
    ORDER BY updated_at
    LIMIT $1;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirp :execrows
DELETE FROM chirps
    WHERE id = $1 AND user_id = $2;
