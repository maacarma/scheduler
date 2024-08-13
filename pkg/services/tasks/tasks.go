package tasks

import (
	"context"
	"fmt"

	models "github.com/maacarma/scheduler/pkg/services/tasks/models"
)

// Repo is the interface that wraps the required repository methods.
// Any underlying database repository should implement these methods.
type Repo interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
	GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error)
	CreateOne(ctx context.Context, task *models.TaskPayload) (string, error)
}

// Service is the interface that wraps tasks service methods.
type Service interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
	GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error)
	Create(ctx context.Context, task *models.TaskPayload) (string, error)
}

type Executor struct {
	task *models.Task
}

func NewExecutor(task *models.Task) *Executor {
	return &Executor{task: task}
}

// tasks is the concrete implementation of the Service interface.
// It holds the required repository instance.
type tasks struct {
	taskRepo Repo
}

// New returns a new instance of the tasks service.
func New(repo Repo) Service {
	return &tasks{repo}
}

func (s *tasks) GetAll(ctx context.Context) ([]*models.Task, error) {
	return s.taskRepo.GetAll(ctx)
}

func (s *tasks) GetByNamespace(ctx context.Context, namespace string) ([]*models.Task, error) {
	return s.taskRepo.GetByNamespace(ctx, namespace)
}

func (s *tasks) Create(ctx context.Context, task *models.TaskPayload) (string, error) {
	if task.Namespace == "" {
		task.Namespace = "default"
	}
	return s.taskRepo.CreateOne(ctx, task)
}

// ExecuteTask executes a task. It returns an error if the task execution fails.
// This method is used by the scheduler to execute tasks.
func (s *Executor) Run() {
	// complete the task execution logic here
	fmt.Println("completed executing task ", s.task.ID)
}
