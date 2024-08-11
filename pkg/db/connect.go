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

type Clients struct {
	Mongo *mongo.Client
	Pg    *pgx.Conn
}

// Connect connects to the database and returns the connection.
func Connect(ctx context.Context, conf *utils.Config) (*Clients, error) {

	db := conf.Database.Db
	pgConnStr := conf.Database.Postgres.ConnString
	mongoConnStr := conf.Database.MongoDB.ConnString
	c := Clients{}

	switch db {
	case "postgres":
		pgConn, err := postgres.Connect(ctx, pgConnStr)
		if err != nil {
			return nil, fmt.Errorf(connErr, db, err)
		}
		c.Pg = pgConn
		return &c, nil

	case "mongodb":
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
