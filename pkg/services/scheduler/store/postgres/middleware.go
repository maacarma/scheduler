package postgres

import (
	"context"
	"encoding/json"
	"fmt"

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
	var t models.Task
	err := json.Unmarshal([]byte(task.Params), &t.Params)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(task.Headers), &t.Headers)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(task.Body), &t.Body)
	if err != nil {
		return nil, err
	}

	return &t, nil
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
	paramsInBytes, err := json.Marshal(task.Headers)
	if err != nil {
		return -1, err
	}

	headersInBytes, err := json.Marshal(task.Headers)
	if err != nil {
		return -1, err
	}

	bodyInBytes, err := json.Marshal(task.Body)
	if err != nil {
		return -1, err
	}

	m := sqlgen.CreateTaskParams{
		Url:       task.Url,
		Method:    task.Method,
		Namespace: task.Namespace,
		Params:    paramsInBytes,
		Headers:   headersInBytes,
		Body:      bodyInBytes,
	}

	fmt.Println("Executing task")

	id, err := r.rawRepo.CreateTask(ctx, m)
	return id, err
}
