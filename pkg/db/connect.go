package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/maacarma/scheduler/config"
	"github.com/maacarma/scheduler/pkg/db/mongodb"
	"github.com/maacarma/scheduler/pkg/db/postgres"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// database names in the config file
	POSTGRES = "postgres"
	MONGO    = "mongo"
	// Error messages
	connErr  = "unable to connect to %s : %v"
	unkDbErr = "unknown database name: %s"
)

type Clients struct {
	Mongo *mongo.Client
	Pg    *pgx.Conn
}

// Connect connects to the database and returns the connection.
func Connect(ctx context.Context, conf *config.Config) (*Clients, error) {

	db := conf.Database.Db
	pgConnStr := conf.Database.Postgres.Url
	mongoConnStr := conf.Database.MongoDB.Url
	c := Clients{}

	switch db {
	case POSTGRES:
		pgConn, err := postgres.Connect(ctx, pgConnStr)
		if err != nil {
			return nil, fmt.Errorf(connErr, db, err)
		}
		c.Pg = pgConn
		return &c, nil

	case MONGO:
		client, err := mongodb.Connect(ctx, mongoConnStr)
		if err != nil {
			return nil, fmt.Errorf(connErr, db, err)
		}
		c.Mongo = client
		return &c, nil

	default:
		return nil, fmt.Errorf(unkDbErr, db)
	}
}
