-- name: InsertUser :one
INSERT INTO users (
    email
) VALUES ($1) RETURNING *;

-- name: FindUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;
