package api

import (
	"context"
	"net/http"
	"time"

	db "github.com/maacarma/scheduler/pkg/db"
	tasks "github.com/maacarma/scheduler/pkg/services/tasks/transport"
	utils "github.com/maacarma/scheduler/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Start starts the API server
func Start(ctx context.Context, logger *zap.Logger, conf *utils.Config) error {

	pgConn, mongoClient, err := db.Connect(ctx, conf)
	if err != nil {
		return err
	}

	defer func() {
		if pgConn != nil {
			pgConn.Close(ctx)
		}
		if mongoClient != nil {
			mongoClient.Disconnect(ctx)
		}
	}()

	r := gin.Default()
	tasks.Activate(r, pgConn, mongoClient)

	errch := make(chan error)
	server := &http.Server{
		Addr:    conf.Application.Port,
		Handler: r,
	}

	go func() {
		logger.Info("Starting server", zap.String("addr", server.Addr))
		errch <- server.ListenAndServe()
	}()

	select {
	case err := <-errch:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
