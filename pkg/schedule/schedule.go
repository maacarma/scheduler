package schedule

import (
	"context"
	"fmt"
	"time"

	db "github.com/maacarma/scheduler/pkg/db"
	svc "github.com/maacarma/scheduler/pkg/services/tasks"
	models "github.com/maacarma/scheduler/pkg/services/tasks/models"
	mongodb "github.com/maacarma/scheduler/pkg/services/tasks/store/mongodb"
	postgres "github.com/maacarma/scheduler/pkg/services/tasks/store/postgres"
	utils "github.com/maacarma/scheduler/utils"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const (
	scheduleErr          = "unable to schedule tasks"
	scheduleSuccess      = "successfully scheduled all tasks"
	scheduledTask        = "successfully scheduled task with id: %s"
	duplicateTask        = "task with id: %s already scheduled"
	deletedTask          = "task with id: %s deleted"
	noTaskFound          = "no task found with id: %s"
	unableToScheduleTask = "unable to schedule task with id: %s due to %v"
)

type repo interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
}

type tasksMap map[string]cron.EntryID

type Scheduler struct {
	repo   repo
	cron   *cron.Cron
	tasks  tasksMap
	conf   *utils.Config
	logger *zap.Logger
}

func New(conf *utils.Config, logger *zap.Logger) (*Scheduler, error) {
	dbClients, err := db.Connect(context.Background(), conf)
	if err != nil {
		return nil, err
	}

	var repo svc.Repo
	switch {
	case dbClients.Pg != nil:
		repo = postgres.New(dbClients.Pg)
	case dbClients.Mongo != nil:
		repo = mongodb.New(dbClients.Mongo)
	}

	cron := cron.New(cron.WithLocation(time.UTC))
	tasks := make(tasksMap)

	return &Scheduler{
		repo:   repo,
		cron:   cron,
		tasks:  tasks,
		conf:   conf,
		logger: logger,
	}, nil
}

func (s *Scheduler) Tasks(ctx context.Context, logger *zap.Logger, conf *utils.Config) error {
	dbClients, err := db.Connect(ctx, conf)
	if err != nil {
		return err
	}

	var repo svc.Repo
	switch {
	case dbClients.Pg != nil:
		repo = postgres.New(dbClients.Pg)
	case dbClients.Mongo != nil:
		repo = mongodb.New(dbClients.Mongo)
	}

	tasks, err := repo.GetAll(ctx)
	for _, t := range tasks {
		if err := s.NewTask(t); err != nil {
			return err
		}
	}

	s.logger.Info(scheduleSuccess)
	return nil
}

func (s *Scheduler) NewTask(t *models.Task) error {
	if _, exists := s.tasks[t.ID]; exists {
		return fmt.Errorf(duplicateTask, t.ID)
	}

	executor := svc.NewExecutor(t)
	entryID, err := s.cron.AddJob(t.Interval, executor)
	if err != nil {
		return fmt.Errorf(unableToScheduleTask, t.ID, err)
	}

	s.logger.Info(fmt.Sprintf(scheduledTask, t.ID))
	s.tasks[t.ID] = entryID
	return nil
}

func (s *Scheduler) DiscardTask(taskID string) error {
	if entryID, exists := s.tasks[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.tasks, taskID)
		s.logger.Info(fmt.Sprintf(deletedTask, taskID))
		return nil
	}

	return fmt.Errorf("task with id: %s not found", taskID)
}
