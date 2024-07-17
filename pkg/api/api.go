package api

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/maacarma/scheduler/pkg/services/scheduler/store/postgres"
	scheduler "github.com/maacarma/scheduler/pkg/services/scheduler/transport"
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
