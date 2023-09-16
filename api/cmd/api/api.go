package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/LeonLow97/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Open Database Connection
	db, err := openDB()
	if err != nil {
		log.Fatalln(err)
	}

	router := routes(db)

	log.Println("Application has started. Listening port is 4000")
	http.ListenAndServe(":4000", router)
}

func openDB() (*sqlx.DB, error) {
	environment := "production" // TODO: make this dynamic, get from env file?
	var databaseURL string

	config, err := config.LoadConfig(environment)
	if err != nil {
		return nil, err
	}

	if environment == "development" {
		databaseURL = config.Development.URL
	} else {
		databaseURL = config.Production.URL
	}

	connConfig, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		errMsg := err.Error()
		errMsg = regexp.MustCompile(`(://[^:]+:).+(@.+)`).ReplaceAllString(errMsg, "$1*****$2")
		errMsg = regexp.MustCompile(`(password=).+(\s+)`).ReplaceAllString(errMsg, "$1*****$2")
		return nil, fmt.Errorf("parsing DSN failed: %s", errMsg)
	}
	connectionStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sqlx.Open("pgx", connectionStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	instance, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"leon-db",
		instance,
	)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return db, nil
}
