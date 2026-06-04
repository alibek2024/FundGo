-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, 
    first_name, 
    last_name
) VALUES (
    $1, $2, $3, $4
)
RETURNING  *;

-- name: UpdateInfo :one
UPDATE users 
SET 
    first_name = $1,
    last_name = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING  *;

-- name: UpdateEmail :one
UPDATE users 
SET 
    email = $1
WHERE id = $2
RETURNING  *;

-- name: UpdatePassword :one
UPDATE users 
SET 
    password_hash = $1
WHERE id = $2
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

-- name: RestoreUser :one
UPDATE users 
SET deleted_at = NULL 
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;