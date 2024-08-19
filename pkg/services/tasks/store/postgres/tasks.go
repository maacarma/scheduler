package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	models "github.com/maacarma/scheduler/pkg/services/tasks/models"
	sqlgen "github.com/maacarma/scheduler/pkg/services/tasks/store/postgres/sqlgen"

	"github.com/jackc/pgx/v5"
)

// repo is the concrete implementation of the Tasks Repo interface.
// It holds the required querier instance, which wraps the sqlgen methods.
type repo struct {
	querier sqlgen.Querier
}

// New returns a new instance of the postgres repo.
func New(pgConn *pgx.Conn) *repo {
	querier := sqlgen.New(pgConn)
	return &repo{querier: querier}
}

// GetAll returns all tasks from the database.
func (r *repo) GetAll(ctx context.Context) ([]*models.Task, error) {
	tasks, err := r.querier.GetTasks(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Task, 0)
	for _, task := range tasks {
		t, err := convert(task)
		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

// GetByNamespace returns all tasks from the database with the given namespace.
func (r *repo) GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error) {
	tasks, err := r.querier.GetTasksByNamespace(ctx, namespace)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Task, 0)
	for _, task := range tasks {
		t, err := convert(task)
		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

// CreateOne creates a new task and returns the id.
func (r *repo) CreateOne(ctx context.Context, task *models.TaskPayload) (string, error) {
	paramsInBytes, err := json.Marshal(task.Params)
	if err != nil {
		return "", err
	}

	headersInBytes, err := json.Marshal(task.Headers)
	if err != nil {
		return "", err
	}

	bodyInBytes, err := json.Marshal(task.Body)
	if err != nil {
		return "", err
	}

	m := sqlgen.CreateTaskParams{
		Url:       task.Url,
		Method:    task.Method,
		Namespace: task.Namespace,
		Params:    paramsInBytes,
		Headers:   headersInBytes,
		Body:      bodyInBytes,
		StartUnix: task.StartUnix,
		EndUnix:   task.EndUnix,
		Interval:  task.Interval,
	}

	id, err := r.querier.CreateTask(ctx, m)
	return fmt.Sprint(id), err
}

// convert converts a sqlgen task to a native task model.
func convert(task *sqlgen.Task) (*models.Task, error) {
	var t models.Task
	err := json.Unmarshal(task.Params, &t.Params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(task.Headers, &t.Headers)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(task.Body, &t.Body)
	if err != nil {
		return nil, err
	}

	t.ID = fmt.Sprint(task.ID)
	t.Url = task.Url
	t.Method = task.Method
	t.Namespace = task.Namespace
	t.StartUnix = task.StartUnix
	t.EndUnix = task.EndUnix
	t.Interval = task.Interval

	return &t, nil
}
