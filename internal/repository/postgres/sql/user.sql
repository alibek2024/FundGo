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

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: AddBalance :exec
UPDATE users 
SET balance = balance + $2
WHERE id = $1;

-- name: SubtractBalance :execrows
UPDATE users 
SET balance = balance - $2
WHERE id = $1
AND balance >= $2;

-- name: GetByID :one
SELECT *
FROM users
WHERE id = $1
AND deleted_at IS NULL;

-- name: GetByEmail :one
SELECT *
FROM users
WHERE email = $1
AND deleted_at IS NULL;

-- name: GetBalance :one
SELECT balance
FROM users
WHERE id = $1
AND deleted_at IS NULL; 

-- name: UserResponce :one 
SELECT id, email, first_name, last_name, balance, created_at, updated_at, deleted_at
FROM users
WHERE id = $1;

-- name: RestoreUser :one
UPDATE users 
SET deleted_at = NULL 
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;