package infrastructure

import (
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresConn struct {
	*sqlx.DB
}

func NewPostgresConn(cfg *Config) (*PostgresConn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.PostgresUser,
		cfg.Postgres.PostgresPassword,
		cfg.Postgres.PostgresHost,
		cfg.Postgres.PostgresPort,
		cfg.Postgres.PostgresDB,
	)

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		errMsg := err.Error()
		errMsg = regexp.MustCompile(`(://[^:]+:).+(@.+)`).ReplaceAllString(errMsg, "$1*****$2")
		errMsg = regexp.MustCompile(`(password=).+(\s+)`).ReplaceAllString(errMsg, "$1*****$2")
		return nil, fmt.Errorf("parsing DSN failed: %s", errMsg)
	}
	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	dbConn := &PostgresConn{db}

	return dbConn, nil
}
