-- name: CreateUser :one
INSERT INTO users (business_id,
                   email,
                   phone,
                   full_name,
                   password_hash,
                   role)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
  AND business_id = $2;

-- name: GetUserByPhone :one
SELECT *
FROM users
WHERE phone = $1
  AND business_id = $2;

-- name: UpdateUser :one
UPDATE users
SET email     = COALESCE($2, email),
    phone     = COALESCE($3, phone),
    full_name = COALESCE($4, full_name)
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2
WHERE id = $1;

-- name: CreateBusinessAdmin :one
WITH created_business AS (
    INSERT INTO businesses (name)
        VALUES ($1)
        RETURNING id)
INSERT
INTO users (business_id,
            email,
            phone,
            full_name,
            password_hash,
            role)
VALUES ((SELECT id FROM created_business),
        $2,
        $3,
        $4,
        $5,
        'admin')
RETURNING *;