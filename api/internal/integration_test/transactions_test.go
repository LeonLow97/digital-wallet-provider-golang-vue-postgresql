package integration_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// db credentials
var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "database_test"
	port     = "5433" // 5432 might be used for actual Postgres server
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=30"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sqlx.DB

func TestMain(m *testing.M) {
	// connect to docker and create a pool to manage Docker containers in tests
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}
	pool = p

	// setting options for container configuration
	options := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		},
	}

	// running docker container with the specified options
	// `resource` object encapsulates information and methods related to the container
	// `resource` object includes details such as container's ID, IP address, and other metadata
	// Can use this object to perform various operations on the container, such as stopping,
	// purging (removing), inspecting  or checking its status.
	resource, err = pool.RunWithOptions(&options)
	if err != nil {
		// purge removes a container and linked volumes from docker
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// `pool.Retry` is being used to repeatedly attempt a specific operation until it
	// succeeds or a certain number of retry attempts is exhausted. In this case, we are using
	// it to wait for the database to become available, before proceeding with the tests
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sqlx.Connect("postgres", fmt.Sprintf(dsn, host, port, user, password, dbName))
		pool.MaxWait = 20 * time.Minute
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		// purge the container if we still cannot connect to postgres database
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	// create db tables
	err = createTables()
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("error creating tables: %s", err)
	}

	// seed database
	err = seedDatabase()
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("error seeding the database: %s", err)
	}

	// this is going to run all the tests for the entire package; ensure that all the setup
	// needed by all package tests id done by this point
	code := m.Run()

	// `os.Exit` is a function that immediately terminates the program and
	// does not allow any deferred function to execute.
	// `pool.Purge` will free up previously used resources and subsequently kill the container(s)
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	createTablesScriptBytes, err := os.ReadFile("../testdata/init.sql")
	if err != nil {
		return err
	}

	_, err = testDB.Exec(string(createTablesScriptBytes))
	if err != nil {
		return err
	}

	return nil
}

func seedDatabase() error {
	seedDatabaseScriptBytes, err := os.ReadFile("../testdata/seed.sql")
	if err != nil {
		return err
	}

	_, err = testDB.Exec(string(seedDatabaseScriptBytes))
	if err != nil {
		return err
	}

	return nil
}

func Test_PingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("unable to ping database")
	}
}
