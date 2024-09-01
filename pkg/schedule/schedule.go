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
	GetActiveTasks(ctx context.Context, curUnix utils.Unix) ([]*models.Task, error)
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
	dbClients, err := db.Connect(ctx, conf)
	if err != nil {
		return nil, err
	}

	var repo repo
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
	tasks, err := s.repo.GetActiveTasks(s.ctx, utils.CurrentUTCUnix())
	if err != nil {
		return fmt.Errorf(scheduleErr, err)
	}

	for _, t := range tasks {
		s.ScheduleTask(t)
	}

	s.cron.Start()
	s.logger.Info(scheduleSuccess)
	return nil
}

// ScheduleTask schedules the task based on the start time.
func (s *Scheduler) ScheduleTask(t *models.Task) {
	curUnix := utils.CurrentUTCUnix()
	startUnix := utils.Unix(t.StartUnix)
	if utils.CurrentUTCUnix() < startUnix {
		go s.scheduleTaskWithDelay(startUnix.Sub(curUnix, false), t)
		return
	}

	s.scheduleExistingTask(t)
}

// ScheduleTaskNow adds the task to the scheduler.
// Runs the task immediately because cron/v3 doesn't support immediate scheduling.
// and also triggers a goroutine to discard the task after the end time.

// It returns an error if the task is already scheduled.
func (s *Scheduler) ScheduleTaskNow(t *models.Task) error {
	endUnix := utils.Unix(t.EndUnix)
	curUnix := utils.CurrentUTCUnix()

	if _, exists := s.tasks[t.ID]; exists {
		return fmt.Errorf(duplicateTask, t.ID)
	}

	executor := svc.NewExecutor(t, s.logger)
	updatedInterval := utils.ConvertToCronInterval(t.Interval)
	// runs the task in separate goroutine, this shouldn't be blocking
	go executor.Run()
	entryID, err := s.cron.AddJob(updatedInterval, executor)
	if err != nil {
		return fmt.Errorf(unableToScheduleTask, t.ID, err)
	}

	s.tasks[t.ID] = entryID
	go s.discardTaskWithDelay(endUnix.Sub(curUnix, false), t.ID)

	return nil
}

// scheduleTaskWithDelay schedules the task after the duration.
// calls ScheduleTaskNow after the duration.
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
			s.logger.Info(fmt.Sprintf(scheduledTask, t.ID))
			return
		}
	}
}

// scheduleExistingTask schedules the existing task.
// It calculates the next trigger time based on the current time and the start time.
//
// beware: panics if the task.StartUnix is greater than the current time.
func (s *Scheduler) scheduleExistingTask(t *models.Task) {
	startUnix := utils.Unix(t.StartUnix)
	endUnix := utils.Unix(t.EndUnix)
	curUnix := utils.CurrentUTCUnix()

	// parsing the interval according to the cron @every format
	interval, _ := time.ParseDuration(t.Interval)
	updatedInterval := cron.Every(interval).Delay
	intervalInSeconds := int64(updatedInterval.Seconds())

	nextTrigger := time.Duration(intervalInSeconds-(int64(curUnix-startUnix)%intervalInSeconds)) * time.Second
	endDuration := endUnix.Sub(curUnix, false)
	if nextTrigger > endDuration {
		return
	}

	go s.scheduleTaskWithDelay(nextTrigger, t)
}

// DiscardTaskNow removes a task from the scheduler
// and won't stops the task if it is running.
//
// if the task is not found in scheduler, it logs a message.
func (s *Scheduler) DiscardTaskNow(taskID string) {
	if entryID, exists := s.tasks[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.tasks, taskID)
		s.logger.Info(fmt.Sprintf(deletedTask, taskID))
		return
	}

	s.logger.Info(fmt.Sprintf(noActiveTaskFoundToDiscard, taskID))
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
