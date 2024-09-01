package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	models "github.com/maacarma/scheduler/pkg/services/tasks/models"
	sqlgen "github.com/maacarma/scheduler/pkg/services/tasks/store/postgres/sqlgen"
	utils "github.com/maacarma/scheduler/utils"

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

// GetByID returns a task from the database with the given id.
func (r *repo) GetByID(ctx context.Context, idStr string) (*models.Task, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	task, err := r.querier.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	t, err := convert(task)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *repo) GetActiveTasks(ctx context.Context, curUnix utils.Unix) ([]*models.Task, error) {
	tasks, err := r.querier.GetActiveTasks(ctx, int64(curUnix))
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
		Paused:    task.Paused,
	}

	id, err := r.querier.CreateTask(ctx, m)
	return fmt.Sprint(id), err
}

// UpdateStatus updates the paused status of a task.
func (r *repo) UpdateStatus(ctx context.Context, idStr string, paused bool) error {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	args := sqlgen.UpdateTaskStatusParams{ID: id, Paused: paused}
	return r.querier.UpdateTaskStatus(ctx, args)
}

// DeleteTask deletes a task
func (r *repo) Delete(ctx context.Context, idStr string) error {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}

	return r.querier.DeleteTask(ctx, id)
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
	t.Paused = task.Paused

	return &t, nil
}
