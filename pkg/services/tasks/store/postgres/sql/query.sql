-- name: GetTasks :many
SELECT * FROM tasks;

-- name: GetTasksByNamespace :many
SELECT * FROM tasks
WHERE namespace = $1;

-- name: CreateTask :one
INSERT INTO tasks (
  url, method, namespace, params, headers, body, start_unix, end_unix, interval
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING _id;
