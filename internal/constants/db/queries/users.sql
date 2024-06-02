-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: LoginUser :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;