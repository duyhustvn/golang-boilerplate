package authsqlrepo_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"boilerplate/internal/config"
	"boilerplate/internal/logger"
	authsqlrepo "boilerplate/internal/modules/auth/repository/sql"
	postgresql "boilerplate/pkg/postgresql"
)

func TestSqlRepo(t *testing.T) {
	ctx := context.Background()

	cfg := config.Config{
		Env: config.Env{
			Environment: "test",
		},
	}

	log, err := logger.GetLogger(&cfg)
	if err != nil {
		t.Fatalf("failed to init logger %+v", err)
	}

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresC, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Run any migrations on the database
	_, _, err = postgresC.Exec(ctx, []string{"psql", "-U", dbUser, "-d", dbName, "-c", "CREATE TABLE users (id SERIAL, username VARCHAR UNIQUE NOT NULL, password VARCHAR NOT NULL)"})
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create a snapshot of the database to restore later
	err = postgresC.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container
	defer func() {
		if err := postgresC.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	dbURL, err := postgresC.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Test inserting a user", func(t *testing.T) {
		t.Cleanup(func() {
			// 3. In each test, reset the DB to its snapshot state.
			err = postgresC.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})

		sqlxDB, err := postgresql.NewSqlx(dbURL, log)
		if err != nil {
			t.Fatal(err)
		}

		repo := authsqlrepo.NewSqlRepo(sqlxDB, log)

		tests := []struct {
			newUsername string
			newPassword string
			inserted    int64
			hasError    bool
		}{
			{"admin", "changeme1", 1, false},
			{"admin", "changeme2", 0, true}, // should got error because inserted the duplicated username
		}

		for _, tt := range tests {
			// Insert new user with same username
			inserted, err := repo.SaveNewUser(ctx, tt.newUsername, tt.newPassword)

			if tt.inserted != inserted {
				t.Fatalf("Expected %d to equal %d", tt.inserted, inserted)
			}

			if (err != nil) != tt.hasError {
				t.Fatalf("Expected %t to equal %t", tt.hasError, err != nil)
			}
		}

		var username string
		var password string
		err = sqlxDB.QueryRowContext(context.Background(), "SELECT username, password FROM users LIMIT 1").Scan(&username, &password)
		if err != nil {
			t.Fatal(err)
		}

		if username != "admin" {
			t.Fatalf("Expected %s to equal `admin`", username)
		}

		if password != "changeme1" {
			t.Fatalf("Expected %s to equal `changme1`", password)
		}
	})

	// 4. Run as many tests as you need, they will each get a clean database
	t.Run("Test querying empty DB", func(t *testing.T) {
		t.Cleanup(func() {
			// 3. In each test, reset the DB to its snapshot state.
			err = postgresC.Restore(ctx)
			if err != nil {
				t.Fatal(err)
			}
		})
		sqlxDB, err := postgresql.NewSqlx(dbURL, log)
		if err != nil {
			t.Fatal(err)
		}

		// repo := authsqlrepo.NewSqlRepo(sqlxDB, log)

		var username string
		var password int64
		err = sqlxDB.QueryRowContext(context.Background(), "SELECT username, password FROM users LIMIT 1").Scan(&username, &password)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("Expected error to be a NoRows error, since the DB should be empty on every test. Got %s instead", err)
		}
	})
}
