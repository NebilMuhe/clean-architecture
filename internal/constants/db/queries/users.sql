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

-- name: CheckUserExists :one
SELECT EXISTS(
SELECT *
FROM users
WHERE username = $1 OR email = $2
LIMIT 1
);



-- name: CreateSession :one
INSERT INTO sessions (
  username,
  refresh_token
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetToken :one
SELECT *
FROM sessions
WHERE username = $1 and deleted_at ISNULL
LIMIT 1;

-- name: DeleteRefreshToken :one
UPDATE sessions 
SET deleted_at = NOW() 
WHERE username = $1
RETURNING *;