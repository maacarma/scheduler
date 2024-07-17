package scheduler

import (
	"context"

	"github.com/maacarma/scheduler/pkg/services/scheduler/store/postgres/sqlgen"
)

type Repo interface {
	ListTasks(ctx context.Context) ([]*sqlgen.Task, error)
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
	return s.repo.ListTasks(ctx)
}
