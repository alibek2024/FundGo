-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, 
    first_name, 
    last_name
) VALUES (
    $1, $2, $3, $4
)
RETURNING  *;

-- name: UpdateUser :one
UPDATE users 
SET 
    email = $1,
    password_hash = $2,
    first_name = $3,
    last_name = $4,
    updated_at = NOW()
WHERE id = $5
RETURNING  *;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;

-- name: TopUp :exec
UPDATE users 
SET balance = balance + $2
WHERE id = $1;

-- name: Withdraw :execrows
UPDATE users 
SET balance = balance - $2
WHERE id = $1
AND balance >= $2;

-- name: GetByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetByEmail :one
SELECT *
FROM users
WHERE email = $1;