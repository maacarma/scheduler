// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sqlgen

import (
	"context"
)

const listTasks = `-- name: ListTasks :many
SELECT id, name, desciption FROM tasks
ORDER BY name
`

func (q *Queries) ListTasks(ctx context.Context) ([]*Task, error) {
	rows, err := q.db.Query(ctx, listTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(&i.ID, &i.Name, &i.Desciption); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
