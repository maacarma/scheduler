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

	scheduler, err := schedule.New(config, logger)
	if err != nil {
		logger.Fatal("Cannot create scheduler", zap.Error(err))
	}
	scheduler.Start(ctx)

	if err := api.Start(ctx, logger, config); err != nil {
		logger.Fatal("Cannot start api server", zap.Error(err))
	}
}
