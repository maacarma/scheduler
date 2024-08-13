package schedule

import (
	"context"
	"time"

	"github.com/maacarma/scheduler/pkg/db"
	svc "github.com/maacarma/scheduler/pkg/services/tasks"
	"github.com/maacarma/scheduler/pkg/services/tasks/store/mongodb"
	"github.com/maacarma/scheduler/pkg/services/tasks/store/postgres"
	"github.com/maacarma/scheduler/utils"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const (
	scheduleErr          = "unable to schedule tasks"
	scheduleSuccess      = "successfully scheduled tasks"
	scheduledTask        = "successfully scheduled task with id: %s"
	unableToScheduleTask = "unable to schedule task with id: %s"
)

// make a new schedule struct

func Tasks(ctx context.Context, logger *zap.Logger, conf *utils.Config) error {

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

	// main cron instance
	c := cron.New(cron.WithLocation(time.UTC))

	tasks, err := repo.GetAll(ctx)
	for _, t := range tasks {
		executor := svc.NewExecutor(t)
		c.AddJob(t.Interval, executor)
	}

	return nil

}
