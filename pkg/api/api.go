package api

import (
	"context"
	"log"

	postgres "github.com/maacarma/scheduler/pkg/db/postgres"
	scheduler "github.com/maacarma/scheduler/pkg/services/scheduler/transport"

	"github.com/gin-gonic/gin"
)

func Start() {

	pgConn, err := postgres.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if pgConn != nil {
			pgConn.Close(context.Background())
		}
	}()

	r := gin.Default()
	scheduler.Activate(r, pgConn)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
