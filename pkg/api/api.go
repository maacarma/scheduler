package api

import (
	"context"
	"net/http"
	"time"

	config "github.com/maacarma/scheduler/config"
	db "github.com/maacarma/scheduler/pkg/db"
	svc "github.com/maacarma/scheduler/pkg/services/tasks"
	tasks "github.com/maacarma/scheduler/pkg/services/tasks/transport"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	closingErr = "unable to close the server"
)

// Start starts the API server
func Start(ctx context.Context, scheduler svc.Scheduler, logger *zap.Logger, conf *config.Config) error {

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
