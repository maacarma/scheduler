package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect connects to the mongodb server and returns the client.
// It checks the connection by pinging the server.
// It returns an error if the connection/ping fails.
func Connect(ctx context.Context, connString string) (*mongo.Client, error) {
	timeout := time.Second * 5
	opts := options.Client().SetServerSelectionTimeout(timeout)

	client, err := mongo.Connect(ctx, opts.ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	pingErr := client.Ping(ctx, nil)
	if pingErr != nil {
		return nil, pingErr
	}

	return client, nil
}
