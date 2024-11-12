-- name: CreateService :one
INSERT INTO services (business_id,
                      name,
                      description,
                      duration,
                      price,
                      is_active)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetService :one
SELECT *
FROM services
WHERE id = $1;

-- name: UpdateService :one
UPDATE services
SET name        = $2,
    description = $3,
    duration    = $4,
    price       = $5,
    is_active   = $6
WHERE id = $1
RETURNING *;

-- name: DeleteService :exec
UPDATE services
SET is_active = false
WHERE id = $1;

-- name: ListServices :many
SELECT *
FROM services
WHERE business_id = $1
  AND is_active = true
ORDER BY name;

-- name: ListActiveServices :many
SELECT *
FROM services
WHERE business_id = $1
  AND is_active = true
ORDER BY name;