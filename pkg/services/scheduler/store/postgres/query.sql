-- name: GetTasks :many
SELECT * FROM tasks;

-- name: GetTasksByNamespace :many
SELECT * FROM tasks
WHERE namespace = $1;

-- name: CreateTask :one
INSERT INTO tasks (
  url, method, namespace, params, headers, body
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id;
