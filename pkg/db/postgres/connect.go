package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// Connect creates a connection to the postgres server.
// It checks the connection by pinging the server.
// It returns an error if the connection/ping fails.
func Connect(ctx context.Context, connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	if err := initialize(ctx, conn); err != nil {
		return nil, fmt.Errorf("error initializing postgres tables: %w", err)
	}

	return conn, nil
}

// initialize creates the schema in the postgres database.
func initialize(ctx context.Context, pgxConn *pgx.Conn) error {
	path := "pkg/services/tasks/store/postgres/sql/schema.sql"
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
