-- Schedule Templates
-- name: CreateTemplate :one
INSERT INTO schedule_templates (employee_id,
                                day_of_week,
                                start_time,
                                end_time,
                                is_break)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateTemplate :one
UPDATE schedule_templates
SET start_time = $2,
    end_time   = $3,
    is_break   = $4
WHERE id = $1
RETURNING *;

-- name: DeleteTemplate :exec
DELETE
FROM schedule_templates
WHERE id = $1;

-- name: ListTemplates :many
SELECT *
FROM schedule_templates
WHERE employee_id = $1
ORDER BY day_of_week, start_time;

-- Schedule Overrides
-- name: CreateOverride :one
INSERT INTO schedule_overrides (employee_id,
                                override_date,
                                start_time,
                                end_time,
                                is_working_day,
                                is_break)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateOverride :one
UPDATE schedule_overrides
SET start_time     = $2,
    end_time       = $3,
    is_working_day = $4,
    is_break       = $5
WHERE id = $1
RETURNING *;

-- name: DeleteOverride :exec
DELETE
FROM schedule_overrides
WHERE id = $1;

-- name: ListOverrides :many
SELECT *
FROM schedule_overrides
WHERE employee_id = $1
  AND override_date BETWEEN $2 AND $3
ORDER BY override_date, start_time;

-- name: GetEmployeeSchedule :one
SELECT *
FROM schedule_templates
WHERE employee_id = $1 AND day_of_week = $2;