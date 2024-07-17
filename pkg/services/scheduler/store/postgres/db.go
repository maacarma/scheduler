package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, "postgres://scheduler:scheduler@localhost:5432/scheduler")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	return conn, nil
}
