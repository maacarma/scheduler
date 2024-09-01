-- name: GetTasks :many
SELECT * FROM tasks;

-- name: GetTasksByNamespace :many
SELECT * FROM tasks
WHERE namespace = $1;

-- name: GetActiveTasks :many
SELECT * FROM tasks
WHERE end_unix >= $1 AND NOT paused;

-- name: CreateTask :one
INSERT INTO tasks (
  url, method, namespace, params, headers, body, start_unix, end_unix, interval, paused
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING _id;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE _id = $1;

-- name: UpdateTaskStatus :exec
UPDATE tasks
SET paused = $2
WHERE _id = $1;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE _id = $1;