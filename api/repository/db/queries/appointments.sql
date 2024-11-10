-- name: CreateAppointment :one
INSERT INTO appointments (business_id,
                          client_id,
                          employee_id,
                          service_id,
                          start_time,
                          end_time,
                          status,
                          reminder_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAppointment :one
SELECT a.*,
       c.email     as client_email,
       c.phone     as client_phone,
       c.full_name as client_full_name,
       e.email     as employee_email,
       e.phone     as employee_phone,
       e.full_name as employee_full_name,
       s.name      as service_name,
       s.duration  as service_duration,
       s.price     as service_price
FROM appointments a
         JOIN users c ON c.id = a.client_id
         JOIN users e ON e.id = (SELECT user_id FROM employees WHERE id = a.employee_id)
         JOIN services s ON s.id = a.service_id
WHERE a.id = $1;

-- name: UpdateAppointment :one
UPDATE appointments
SET start_time    = $2,
    end_time      = $3,
    status        = $4,
    reminder_time = $5
WHERE id = $1
RETURNING *;

-- name: CancelAppointment :one
UPDATE appointments
SET status = 'cancelled'
WHERE id = $1
RETURNING *;

-- name: ListBusinessAppointments :many
SELECT a.*,
       c.email     as client_email,
       c.phone     as client_phone,
       c.full_name as client_full_name,
       e.email     as employee_email,
       e.phone     as employee_phone,
       e.full_name as employee_full_name,
       s.name      as service_name,
       s.duration  as service_duration,
       s.price     as service_price
FROM appointments a
         JOIN users c ON c.id = a.client_id
         JOIN users e ON e.id = (SELECT user_id FROM employees WHERE id = a.employee_id)
         JOIN services s ON s.id = a.service_id
WHERE a.business_id = $1
  AND a.start_time BETWEEN $2 AND $3
ORDER BY a.start_time;

-- name: ListEmployeeAppointments :many
SELECT a.*,
       c.email     as client_email,
       c.phone     as client_phone,
       c.full_name as client_full_name,
       e.email     as employee_email,
       e.phone     as employee_phone,
       e.full_name as employee_full_name,
       s.name      as service_name,
       s.duration  as service_duration,
       s.price     as service_price
FROM appointments a
         JOIN users c ON c.id = a.client_id
         JOIN users e ON e.id = (SELECT user_id FROM employees WHERE id = a.employee_id)
         JOIN services s ON s.id = a.service_id
WHERE a.employee_id = $1
  AND a.start_time BETWEEN $2 AND $3
ORDER BY a.start_time;

-- name: ListClientAppointments :many
SELECT a.*,
       c.email     as client_email,
       c.phone     as client_phone,
       c.full_name as client_full_name,
       e.email     as employee_email,
       e.phone     as employee_phone,
       e.full_name as employee_full_name,
       s.name      as service_name,
       s.duration  as service_duration,
       s.price     as service_price
FROM appointments a
         JOIN users c ON c.id = a.client_id
         JOIN users e ON e.id = (SELECT user_id FROM employees WHERE id = a.employee_id)
         JOIN services s ON s.id = a.service_id
WHERE a.client_id = $1
  AND a.start_time BETWEEN $2 AND $3
ORDER BY a.start_time;

-- name: CheckEmployeeAvailability :one
SELECT COUNT(*) = 0 as is_available
FROM appointments
WHERE employee_id = $1
  AND status = 'scheduled'
  AND (
    (start_time, end_time) OVERLAPS ($2, $3)
        OR
    (start_time, end_time) OVERLAPS ($2, $3)
    );