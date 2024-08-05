package postgres

import (
	"context"
	"encoding/json"

	models "github.com/maacarma/scheduler/pkg/services/scheduler/models"
	sqlgen "github.com/maacarma/scheduler/pkg/services/scheduler/store/postgres/sqlgen"

	"github.com/jackc/pgx/v5"
)

type RawRepository interface {
	GetTasks(ctx context.Context) ([]*sqlgen.Task, error)
	GetTasksByNamespace(ctx context.Context, namespace string) ([]*sqlgen.Task, error)
	CreateTask(ctx context.Context, arg sqlgen.CreateTaskParams) (int64, error)
}

type Repository struct {
	rawRepo RawRepository
}

func New(pgConn *pgx.Conn) *Repository {
	sql := sqlgen.New(pgConn)
	return &Repository{rawRepo: sql}
}

func convertToTaskModel(task *sqlgen.Task) (*models.Task, error) {
	t := models.Task{
		ID:        task.ID,
		Url:       task.Url,
		Method:    task.Method,
		Namespace: task.Namespace,
	}

	err := json.Unmarshal([]byte(task.Headers), &t.Headers)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(task.Body), &t.Body)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func convertToDBModel(task *models.TaskPayload) (*sqlgen.CreateTaskParams, error) {
	var m sqlgen.CreateTaskParams

	headersInBytes, err := json.Marshal(task.Headers)
	if err != nil {
		return nil, err
	}

	bodyInBytes, err := json.Marshal(task.Body)
	if err != nil {
		return nil, err
	}

	m = sqlgen.CreateTaskParams{
		Url:       task.Url,
		Method:    task.Method,
		Namespace: task.Namespace,
		Headers:   headersInBytes,
		Body:      bodyInBytes,
	}

	return &m, nil
}

func (r *Repository) GetTasks(ctx context.Context) ([]*models.Task, error) {
	tasks, err := r.rawRepo.GetTasks(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Task, 0)
	for _, task := range tasks {
		t, err := convertToTaskModel(task)
		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (r *Repository) GetTasksByNamespace(ctx context.Context, namespace string) ([]*models.Task, error) {
	tasks, err := r.rawRepo.GetTasksByNamespace(ctx, namespace)
	if err != nil {
		return nil, err
	}

	result := make([]*models.Task, 0)
	for _, task := range tasks {
		t, err := convertToTaskModel(task)
		if err != nil {
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (r *Repository) CreateTask(ctx context.Context, task *models.TaskPayload) (int64, error) {
	m, err := convertToDBModel(task)
	if err != nil {
		return 0, err
	}

	id, err := r.rawRepo.CreateTask(ctx, *m)
	return id, err
}
