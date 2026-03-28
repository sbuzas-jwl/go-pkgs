// sqlite is a facade over the data storage layer.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sbuzas-jwl/go-pkgs/todo/pkg/logging"
)

type DB struct {
	Pool *sqlx.DB
}

// NewFromEnv sets up the database connections using the configuration in the
// process's environment variables. This should be called just once per server
// instance.
func NewFromEnv(ctx context.Context, cfg *Config) (*DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbDSN(cfg).String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}
	dbx := sqlx.NewDb(sqlDB, "sqlite3")

	return &DB{Pool: dbx}, nil
}

// Close releases database connections.
func (db *DB) Close(ctx context.Context) {
	logger := logging.FromContext(ctx)
	logger.Info("Closing connection pool.")
	db.Pool.Close()
}

// dbDSN builds a connection string suitable for the sqlite3 driver, using
// the values of vars.
func dbDSN(cfg *Config) *url.URL {
	return cfg.ConnectionURL()
}
