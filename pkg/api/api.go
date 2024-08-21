package api

import (
	"context"
	"net/http"
	"time"

	db "github.com/maacarma/scheduler/pkg/db"
	"github.com/maacarma/scheduler/pkg/schedule"
	tasks "github.com/maacarma/scheduler/pkg/services/tasks/transport"
	utils "github.com/maacarma/scheduler/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	closingErr = "unable to close the server"
)

// Start starts the API server
func Start(ctx context.Context, scheduler *schedule.Scheduler, logger *zap.Logger, conf *utils.Config) error {

	dbClients, err := db.Connect(ctx, conf)
	if err != nil {
		return err
	}

	r := gin.Default()
	tasks.Activate(r, dbClients, scheduler)

	errch := make(chan error)
	server := &http.Server{
		Addr:    conf.Application.Port,
		Handler: r,
	}

	defer func() {
		logger.Warn("graceful shutting server and db connections")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if dbClients.Pg != nil {
			dbClients.Pg.Close(shutdownCtx)
		}
		if dbClients.Mongo != nil {
			dbClients.Mongo.Disconnect(shutdownCtx)
		}
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Error(closingErr, zap.Error(err))
		}
	}()

	go func() {
		logger.Info("Starting server", zap.String("addr", server.Addr))
		errch <- server.ListenAndServe()
	}()

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		return nil
	}
}
