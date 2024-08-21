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
	scheduleErr                = "unable to schedule tasks due to %v"
	scheduleSuccess            = "successfully scheduled all tasks"
	scheduledTask              = "successfully scheduled task with id: %s"
	duplicateTask              = "task with id: %s already scheduled"
	schedullingAnInactiveTask  = "task with id: %s is inactive"
	deletedTask                = "task with id: %s deleted"
	noTaskFound                = "no task found with id: %s"
	noActiveTaskFoundToDiscard = "no active task found with id: %s to discard"
	unableToScheduleTask       = "unable to schedule task with id: %s due to %v"
)

type repo interface {
	GetAll(ctx context.Context) ([]*models.Task, error)
}

type tasksMap map[string]cron.EntryID

type Scheduler struct {
	ctx    context.Context
	repo   repo
	cron   *cron.Cron
	tasks  tasksMap
	conf   *utils.Config
	logger *zap.Logger
}

// New creates a new scheduler instance.
func New(ctx context.Context, conf *utils.Config, logger *zap.Logger) (*Scheduler, error) {
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
		ctx:    ctx,
		repo:   repo,
		cron:   cron,
		tasks:  tasks,
		conf:   conf,
		logger: logger,
	}, nil
}

// Start starts the scheduler.
// It schedules all the active tasks read from the database.
func (s *Scheduler) Start() error {
	tasks, err := s.repo.GetAll(s.ctx)
	if err != nil {
		return fmt.Errorf(scheduleErr, err)
	}

	for _, t := range tasks {
		err := s.ScheduleTask(t)
		if err != nil {
			return err
		}
	}

	s.cron.Start()
	s.logger.Info(scheduleSuccess)
	return nil
}

// ScheduleTask schedules the new task based on the start time.
func (s *Scheduler) ScheduleTask(t *models.Task) error {
	if utils.CurrentUTCUnix() < t.StartUnix {
		go s.scheduleTaskWithDelay(utils.UTCUnixTimeDiff(t.StartUnix, false), t)
		return nil
	}

	if err := s.ScheduleTaskNow(t); err != nil {
		return err
	}
	return nil
}

// ScheduleTaskNow adds a new task to the scheduler.
// It returns an error if the task is not active or if the task is already scheduled.
// and also triggers a goroutine to discard the task after the end time.
func (s *Scheduler) ScheduleTaskNow(t *models.Task) error {
	if !t.IsActive() {
		s.logger.Warn(fmt.Sprintf(schedullingAnInactiveTask, t.ID))
		return nil
	}

	if _, exists := s.tasks[t.ID]; exists {
		return fmt.Errorf(duplicateTask, t.ID)
	}

	executor := svc.NewExecutor(t)
	updatedInterval := utils.ConvertToCronInterval(t.Interval)
	entryID, err := s.cron.AddJob(updatedInterval, executor)
	if err != nil {
		return fmt.Errorf(unableToScheduleTask, t.ID, err)
	}

	go s.discardTaskWithDelay(utils.UTCUnixTimeDiff(t.EndUnix, false), t.ID)
	s.logger.Info(fmt.Sprintf(scheduledTask, t.ID))
	s.tasks[t.ID] = entryID

	return nil
}

// scheduleTaskWithDelay schedules the task after the duration.
func (s *Scheduler) scheduleTaskWithDelay(duration time.Duration, t *models.Task) {
	ticker := time.NewTicker(duration)

	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			err := s.ScheduleTaskNow(t)
			if err != nil {
				s.logger.Error(fmt.Sprintf(unableToScheduleTask, t.ID, err))
			}
			return
		}
	}
}

// DiscardTaskNow removes a task from the scheduler.
// if the task is not found in scheduler, it logs a message.
func (s *Scheduler) DiscardTaskNow(taskID string) {
	if entryID, exists := s.tasks[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.tasks, taskID)
		s.logger.Info(fmt.Sprintf(deletedTask, taskID))
		return
	}

	s.logger.Warn(fmt.Sprintf(noActiveTaskFoundToDiscard, taskID))
}

// discardAfterEnd discards the task after the duration.
func (s *Scheduler) discardTaskWithDelay(duration time.Duration, taskID string) {
	ticker := time.NewTicker(duration)

	defer ticker.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.DiscardTaskNow(taskID)
			return
		}
	}
}
