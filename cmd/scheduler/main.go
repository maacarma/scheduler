package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/maacarma/scheduler/pkg/api"
	"github.com/maacarma/scheduler/pkg/schedule"
	"github.com/maacarma/scheduler/utils"
	"go.uber.org/zap"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := utils.CreateLogger()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	scheduler, err := schedule.New(ctx, config, logger)
	if err != nil {
		logger.Fatal("unable to create scheduler", zap.Error(err))
	}
	err = scheduler.Start()
	if err != nil {
		logger.Fatal("unable to start scheduler", zap.Error(err))
	}

	if err := api.Start(ctx, scheduler, logger, config); err != nil {
		logger.Fatal("Cannot start api server", zap.Error(err))
	}
}
