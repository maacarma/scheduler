package tasks

import (
	"context"
	"fmt"
	"net/http"

	models "github.com/maacarma/scheduler/pkg/services/tasks/models"
	"github.com/maacarma/scheduler/utils"
)

// Repo is the interface that wraps the required repository methods.
// Any underlying database repository should implement these methods.
type Repo interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
	GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error)
	CreateOne(ctx context.Context, task *models.TaskPayload) (string, error)
}

// Scheduler is the interface that wraps the scheduler methods.
type Scheduler interface {
	ScheduleTask(task *models.Task)
}

// Service is the interface that wraps tasks service methods.
type Service interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
	GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error)
	Create(ctx context.Context, task *models.TaskPayload) (string, int, error)
}

type Executor struct {
	task *models.Task
}

func NewExecutor(task *models.Task) *Executor {
	return &Executor{task: task}
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

	startUnix := utils.Unix(task.StartUnix)
	if startUnix < utils.CurrentUTCUnix() {
		return "", http.StatusBadRequest, fmt.Errorf("start time cannot be in the past")
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

// ExecuteTask executes a task. It returns an error if the task execution fails.
// This method is used by the scheduler to execute tasks.
func (s *Executor) Run() {
	// complete the task execution logic here
	fmt.Println("dry run: executing task ", s.task.ID)
}
