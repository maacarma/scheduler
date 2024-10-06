package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	models "github.com/maacarma/scheduler/pkg/services/tasks/models"
	utils "github.com/maacarma/scheduler/utils"

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

// ExecuteTask executes a task
// This method is used by the cron to execute tasks.
func (s *Executor) Run() {
	s.logger.Info("executing: ", zap.String("task_id", s.task.ID))

	url, err := url.Parse(s.task.Url)
	if err != nil {
		s.logger.Error("failed to parse url", zap.Error(err))
		return
	}
	utils.AppendQueryParams(url, s.task.Params)

	bodyBytes, err := json.Marshal(s.task.Body)
	if err != nil {
		s.logger.Error("failed to marshal body", zap.Error(err))
		return
	}

	req, _ := http.NewRequest(s.task.Method, url.String(), bytes.NewBuffer(bodyBytes))
	req.Header = s.task.Headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("failed to execute task", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	s.logger.Info("task executed", zap.String("task_id", s.task.ID), zap.Int("status_code", resp.StatusCode))
}
