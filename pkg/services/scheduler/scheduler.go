package scheduler

import (
	"context"

	"github.com/maacarma/scheduler/pkg/services/scheduler/store/postgres/sqlgen"
)

type Repo interface {
	GetTasks(ctx context.Context) ([]*sqlgen.Task, error)
	GetTasksByNamespace(ctx context.Context, namespace string) ([]*sqlgen.Task, error)
}

type Service interface {
	GetAll(ctx context.Context) ([]*sqlgen.Task, error)
}

type scheduler struct {
	repo Repo
}

// New Service instance
func New(repo Repo) Service {
	return &scheduler{repo}
}

func (s *scheduler) GetAll(ctx context.Context) ([]*sqlgen.Task, error) {
	return s.repo.GetTasks(ctx)
}
