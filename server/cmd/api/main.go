package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	handlers "github.com/LeonLow97/go-clean-architecture/delivery/http/handler"
	"github.com/LeonLow97/go-clean-architecture/delivery/http/middleware"
	"github.com/LeonLow97/go-clean-architecture/repository"
	"github.com/LeonLow97/go-clean-architecture/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	dbConn, err := connectToDB()
	if err != nil {
		log.Fatalln("error connecting to db", err)
	}
	defer dbConn.Close()

	r := mux.NewRouter()

	r.Use(middleware.CorsMiddleware)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy!"))
	}).Methods(http.MethodGet)

	userRepo := repository.NewUserRepository(dbConn)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	handlers.NewAuthHandler(r, authUsecase)

	port := os.Getenv("SERVICE_PORT")
	log.Println("Server is running on port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Failed to listen to server with error: %v\n", err)
	}
}

// connectToDB connects to PostgreSQL
func connectToDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
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

	return db, nil
}
