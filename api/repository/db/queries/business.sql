-- name: CreateBusiness :one
INSERT INTO businesses (
    name,
    logo_url,
    color_scheme
) VALUES (
             $1,
             $2,
             $3
         ) RETURNING *;

-- name: GetBusiness :one
SELECT * FROM businesses
WHERE id = $1;

-- name: UpdateBusiness :one
UPDATE businesses
SET name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateBusinessAppearance :one
UPDATE businesses
SET logo_url = $2,
    color_scheme = $3
WHERE id = $1
RETURNING *;