-- name: CreateEmployee :one
INSERT INTO employees (business_id,
                       user_id,
                       specialization,
                       is_active)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetEmployee :one
SELECT e.id,
       e.business_id,
       e.user_id,
       e.specialization,
       e.is_active,
       e.created_at,
       u.email      as user_email,
       u.phone      as user_phone,
       u.full_name  as user_full_name,
       u.role       as user_role,
       u.created_at as user_created_at
FROM employees e
         JOIN users u ON u.id = e.user_id
WHERE e.id = $1;

-- name: UpdateEmployee :one
UPDATE employees
SET specialization = COALESCE($2, specialization),
    is_active      = $3
WHERE id = $1
RETURNING *;

-- name: ListEmployees :many
SELECT e.id,
       e.business_id,
       e.user_id,
       e.specialization,
       e.is_active,
       e.created_at,
       u.email      as user_email,
       u.phone      as user_phone,
       u.full_name  as user_full_name,
       u.role       as user_role,
       u.created_at as user_created_at
FROM employees e
         JOIN users u ON u.id = e.user_id
WHERE e.business_id = $1
ORDER BY e.created_at DESC;
-- name: AssignServices :exec
INSERT INTO employee_services (employee_id, service_id)
SELECT $1, unnest($2::int[]);

-- name: RemoveServices :exec
DELETE
FROM employee_services
WHERE employee_id = $1
  AND service_id = ANY ($2::int[]);

-- name: GetEmployeeServices :many
SELECT s.*
FROM services s
         JOIN employee_services es ON es.service_id = s.id
WHERE es.employee_id = $1
ORDER BY s.name;