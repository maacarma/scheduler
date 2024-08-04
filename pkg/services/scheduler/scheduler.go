package scheduler

import (
	"context"

	models "github.com/maacarma/scheduler/pkg/services/scheduler/models"
)

type Repo interface {
	GetTasks(ctx context.Context) ([]*models.Task, error)
	GetTasksByNamespace(ctx context.Context, namespace string) ([]*models.Task, error)
	CreateTask(ctx context.Context, task *models.TaskPayload) (int64, error)
}

type Service interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
	GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error)
	CreateTask(ctx context.Context, task *models.TaskPayload) (int64, error)
}

type scheduler struct {
	repo Repo
}

func New(repo Repo) Service {
	return &scheduler{repo}
}

func (s *scheduler) GetAll(ctx context.Context) ([]*models.Task, error) {
	return s.repo.GetTasks(ctx)
}

func (s *scheduler) GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error) {
	return s.repo.GetTasksByNamespace(ctx, namespace)
}

func (s *scheduler) CreateTask(ctx context.Context, task *models.TaskPayload) (int64, error) {
	return s.repo.CreateTask(ctx, task)
}
