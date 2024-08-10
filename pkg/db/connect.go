package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/maacarma/scheduler/pkg/db/mongodb"
	"github.com/maacarma/scheduler/pkg/db/postgres"
	"github.com/maacarma/scheduler/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// Error messages.
const (
	connErr  = "unable to connect to %s : %v"
	unkDbErr = "unknown database %s"
)

// Connect connects to the database and returns the connection.
func Connect(ctx context.Context, conf *utils.Config) (*pgx.Conn, *mongo.Client, error) {

	db := conf.Database.Db
	pgConnStr := conf.Database.Postgres.ConnString
	mongoConnStr := conf.Database.MongoDB.ConnString

	switch db {
	case "postgres":
		pgConn, err := postgres.Connect(ctx, pgConnStr)
		if err != nil {
			return nil, nil, fmt.Errorf(connErr, db, err)
		}
		return pgConn, nil, nil

	case "mongodb":
		client, err := mongodb.Connect(ctx, mongoConnStr)
		if err != nil {
			return nil, nil, fmt.Errorf(connErr, db, err)
		}
		return nil, client, nil

	default:
		return nil, nil, fmt.Errorf(unkDbErr, db)
	}
}
