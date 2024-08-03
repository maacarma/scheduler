package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, "postgres://postgres:root@localhost:5432/scheduler")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	err = initialize(ctx, conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initialize(ctx context.Context, pgxConn *pgx.Conn) error {
	path := "pkg/services/scheduler/store/postgres/schema.sql"
	c, ioErr := os.ReadFile(path)
	if ioErr != nil {
		return fmt.Errorf("error reading sql file %w", ioErr)
	}

	_, err := pgxConn.Exec(ctx, string(c))
	if err != nil {
		return fmt.Errorf("error executing sql file %w", err)
	}
	return nil
}
