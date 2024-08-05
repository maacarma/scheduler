// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sqlgen

import (
	"context"
)

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (
  url, method, namespace, headers, body
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id
`

type CreateTaskParams struct {
	Url       string `json:"url"`
	Method    string `json:"method"`
	Namespace string `json:"namespace"`
	Headers   []byte `json:"headers"`
	Body      []byte `json:"body"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (int64, error) {
	row := q.db.QueryRow(ctx, createTask,
		arg.Url,
		arg.Method,
		arg.Namespace,
		arg.Headers,
		arg.Body,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getTasks = `-- name: GetTasks :many
SELECT id, url, method, namespace, headers, body FROM tasks
`

func (q *Queries) GetTasks(ctx context.Context) ([]*Task, error) {
	rows, err := q.db.Query(ctx, getTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Task{}
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.Method,
			&i.Namespace,
			&i.Headers,
			&i.Body,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTasksByNamespace = `-- name: GetTasksByNamespace :many
SELECT id, url, method, namespace, headers, body FROM tasks
WHERE namespace = $1
`

func (q *Queries) GetTasksByNamespace(ctx context.Context, namespace string) ([]*Task, error) {
	rows, err := q.db.Query(ctx, getTasksByNamespace, namespace)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Task{}
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.Method,
			&i.Namespace,
			&i.Headers,
			&i.Body,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
