package tasks

import (
	"context"
	"net/http"

	models "github.com/maacarma/scheduler/pkg/services/tasks/models"

	"go.uber.org/zap"
)

// Repo is the interface that wraps the required repository methods.
// Any underlying database repository should implement these methods.
type Repo interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
	GetByID(ctx context.Context, id string) (*models.Task, error)
	GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error)
	CreateOne(ctx context.Context, task *models.TaskPayload) (string, error)
	UpdateStatus(ctx context.Context, id string, paused bool) error
	Delete(ctx context.Context, id string) error
}

// Scheduler is the interface that wraps the scheduler methods.
type Scheduler interface {
	ScheduleTask(task *models.Task)
	DiscardTaskNow(id string)
}

// Service is the interface that wraps tasks service methods.
type Service interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
	GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error)
	Create(ctx context.Context, task *models.TaskPayload) (string, int, error)
	ToggleStatus(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type Executor struct {
	task   *models.Task
	logger *zap.Logger
}

func NewExecutor(task *models.Task, logger *zap.Logger) *Executor {
	return &Executor{task: task, logger: logger}
}

// tasks is the concrete implementation of the Service interface.
// It holds the required repository instance.
type svc struct {
	repo      Repo
	scheduler Scheduler
}

// New returns a new instance of the tasks service.
func New(repo Repo, scheduler Scheduler) Service {
	return &svc{
		repo:      repo,
		scheduler: scheduler,
	}
}

func (s *svc) GetAll(ctx context.Context) ([]*models.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *svc) GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error) {
	return s.repo.GetByNamespace(ctx, namespace)
}

func (s *svc) Create(ctx context.Context, task *models.TaskPayload) (string, int, error) {
	if task.Namespace == "" {
		task.Namespace = "default"
	}

	id, err := s.repo.CreateOne(ctx, task)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	tModel := task.ConvertToTask(id)
	if !tModel.Paused {
		s.scheduler.ScheduleTask(&tModel)
	}

	return id, http.StatusCreated, nil
}

func (s *svc) ToggleStatus(ctx context.Context, id string) error {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	task.Paused = !task.Paused
	err = s.repo.UpdateStatus(ctx, id, task.Paused)
	if err != nil {
		return err
	}

	if task.Paused {
		s.scheduler.DiscardTaskNow(id)
	} else {
		s.scheduler.ScheduleTask(task)
	}

	return nil
}

func (s *svc) Delete(ctx context.Context, id string) error {
	s.scheduler.DiscardTaskNow(id)
	return s.repo.Delete(ctx, id)
}

// ExecuteTask executes a task. It returns an error if the task execution fails.
// This method is used by the scheduler to execute tasks.
func (s *Executor) Run() {
	// complete the task execution logic here
	s.logger.Info("sample run", zap.String("task_id", s.task.ID))
}
