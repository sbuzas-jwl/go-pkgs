package sqlite

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/sethvargo/go-retry"

	// imported to register the sqlite migration driver.
	_ "github.com/mattn/go-sqlite3"
	// imported to register the "file" source migration driver.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	// databaseName is the name of the database.
	databaseName = "test-db"
	// databaseUser and databasePassword are the username and password for
	// connecting to the database. These values are only used for testing.
	databaseUser     = "test-user"
	databasePassword = "testing123"

	// defaultPostgresImageRef is the default database container to use if none is
	// specified.
	defaultPostgresImageRef = "postgres:13-alpine"
)

// ApproxTime is a compare helper for clock skew.
var ApproxTime = cmp.Options{cmpopts.EquateApproxTime(1 * time.Second)}

// TestInstance is a wrapper around the Docker-based database instance.
type TestInstance struct {
	dsnUrl *url.URL

	conn     *sqlx.DB
	connLock sync.Mutex

	skipReason string
}

// MustTestInstance is NewTestInstance, except it prints errors to stderr and
// calls os.Exit when finished. Callers can call Close or MustClose().
func MustTestInstance() *TestInstance {
	testDatabaseInstance, err := NewTestInstance()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return testDatabaseInstance
}

// NewTestInstance creates a new Docker-based database instance. It also creates
// an initial database, runs the migrations, and sets that database as a
// template to be cloned by future tests.
//
// This should not be used outside of testing, but it is exposed in the package
// so it can be shared with other packages. It should be called and instantiated
// in TestMain.
//
// All database tests can be skipped by running `go test -short` or by setting
// the `SKIP_DATABASE_TESTS` environment variable.
func NewTestInstance() (*TestInstance, error) {
	// Querying for -short requires flags to be parsed.
	if !flag.Parsed() {
		flag.Parse()
	}

	// Do not create an instance in -short mode.
	if testing.Short() {
		return &TestInstance{
			skipReason: "🚧 Skipping database tests (-short flag provided)!",
		}, nil
	}

	// Do not create an instance if database tests are explicitly skipped.
	if skip, _ := strconv.ParseBool(os.Getenv("SKIP_DATABASE_TESTS")); skip {
		return &TestInstance{
			skipReason: "🚧 Skipping database tests (SKIP_DATABASE_TESTS is set)!",
		}, nil
	}

	ctx := context.Background()

	cfg := Config{
		Path:     ":memory:",
		User:     databaseUser,
		Password: databasePassword,
	}
	// Build the connection URL.
	connectionURL := cfg.ConnectionURL()

	// Create retryable.
	b := retry.WithMaxRetries(5, retry.NewConstant(1*time.Second))

	// Try to establish a connection to the database, with retries.
	var conn *sqlx.DB
	if err := retry.Do(ctx, b, func(ctx context.Context) error {
		var err error
		conn, err = sqlx.ConnectContext(ctx, "sqlite", connectionURL.String())
		if err != nil {
			return retry.RetryableError(err)
		}
		if err := conn.PingContext(ctx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed waiting for database container to be ready: %w", err)
	}

	// TODO: Run the migrations.
	// if err := dbMigrate(connectionUR); err != nil {
	// 	return nil, fmt.Errorf("failed to migrate database: %w", err)
	// }

	// Return the instance.
	return &TestInstance{
		conn:   conn,
		dsnUrl: connectionURL,
	}, nil
}

// MustClose is like Close except it prints the error to stderr and calls os.Exit.
func (i *TestInstance) MustClose() error {
	if err := i.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return nil
}

// Close terminates the test database instance, cleaning up any resources.
func (i *TestInstance) Close() (retErr error) {
	// Do not attempt to close  things when there's nothing to close.
	if i.skipReason != "" {
		return
	}

	//TODO: implement

	return
}

// NewDatabase creates a new database suitable for use in testing. It returns an
// established database connection and the configuration.
func (i *TestInstance) NewDatabase(tb testing.TB) (*DB, *Config) {
	tb.Helper()

	// Ensure we should actually create the database.
	if i.skipReason != "" {
		tb.Skip(i.skipReason)
	}
	//TODO: implement
	return nil, nil
}
