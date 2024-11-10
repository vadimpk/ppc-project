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

-- name: GetEmployeeSchedule :many
WITH all_dates AS (SELECT generate_series($2::date, $3::date, '1 day'::interval)::date AS date),
     templates AS (SELECT t.id,
                          t.employee_id,
                          t.start_time,
                          t.end_time,
                          t.is_break,
                          EXTRACT(DOW FROM d.date) as day_of_week,
                          d.date
                   FROM all_dates d
                            JOIN schedule_templates t ON EXTRACT(DOW FROM d.date) = t.day_of_week
                   WHERE t.employee_id = $1),
     schedule AS (SELECT date,
                         ARRAY_AGG(
                                 jsonb_build_object(
                                         'id', id,
                                         'start_time', start_time,
                                         'end_time', end_time,
                                         'is_break', is_break,
                                         'is_override', false
                                 )
                                 ORDER BY start_time
                         ) as slots
                  FROM templates
                  GROUP BY date

                  UNION

                  SELECT override_date as date,
                         ARRAY_AGG(
                                 jsonb_build_object(
                                         'id', id,
                                         'start_time', start_time,
                                         'end_time', end_time,
                                         'is_break', is_break,
                                         'is_override', true,
                                         'is_working_day', is_working_day
                                 )
                                 ORDER BY start_time
                         )             as slots
                  FROM schedule_overrides
                  WHERE employee_id = $1
                    AND override_date BETWEEN $2 AND $3
                  GROUP BY override_date)
SELECT *
FROM schedule
ORDER BY date;